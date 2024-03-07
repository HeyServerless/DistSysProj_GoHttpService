package handler

import (
	// "encoding/json"

	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	models "github.com/goService/server/model"
	services "github.com/goService/server/services"
)

func HandleInputExpression(c *gin.Context) {

	log.Println("Expression: ", c.PostForm("expression"))
	ExpressionRequest := models.ExpressionRequest{}

	if err := c.ShouldBindJSON(&ExpressionRequest); err != nil {
		log.Println("Error: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("ExpressionRequest: ", ExpressionRequest)
	expression := ExpressionRequest.Expression

	uuid, _ := services.InsertRequestDataIntoDb(c, &models.CreateRequestRecordModel{RequestBody: expression})
	log.Println("uuid: ", uuid)

	sqsResponse, sqsErr := services.EnqueueRequestToInboundSqs(c, expression, uuid)
	if sqsErr != nil {
		log.Println("Error: ", sqsErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": sqsErr.Error()})
		return
	}
	log.Println("sqsResponse: ", sqsResponse)

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"requestId": uuid,
	})
}

// helper function to check if a string is a number
func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// helper function to check if a character is an operator
func isOperator(c string) bool {
	return c == "+" || c == "-" || c == "*" || c == "/"
}

func breakExpression(expression string) []models.Expression {
	// Divide the expression into smaller sub-expressions based on BODMAS rule
	subExpressions := []string{}
	tmpExp := ""
	for _, c := range expression {
		if c == '(' {
			if len(tmpExp) > 0 {
				subExpressions = append(subExpressions, tmpExp)
				tmpExp = ""
			}
			subExpressions = append(subExpressions, "(")
		} else if c == ')' {
			if len(tmpExp) > 0 {
				subExpressions = append(subExpressions, tmpExp)
				tmpExp = ""
			}
			subExpressions = append(subExpressions, ")")
		} else if c == '+' || c == '-' {
			if len(tmpExp) > 0 {
				subExpressions = append(subExpressions, tmpExp)
				tmpExp = ""
			}
			subExpressions = append(subExpressions, string(c))
		} else if c == '*' || c == '/' {
			if len(tmpExp) > 0 && !strings.Contains(tmpExp, "+") && !strings.Contains(tmpExp, "-") {
				subExpressions = append(subExpressions, tmpExp)
				tmpExp = ""
			}
			tmpExp += string(c)
		} else if c == '^' {
			if len(tmpExp) > 0 && !strings.Contains(tmpExp, "+") && !strings.Contains(tmpExp, "-") && !strings.Contains(tmpExp, "*") && !strings.Contains(tmpExp, "/") {
				subExpressions = append(subExpressions, tmpExp)
				tmpExp = ""
			}
			tmpExp += string(c)
		} else {
			tmpExp += string(c)
		}
	}
	if len(tmpExp) > 0 {
		subExpressions = append(subExpressions, tmpExp)
	}

	// Create independent expressions with unique IDs
	independentExpressions := []models.Expression{}
	for i, exp := range subExpressions {
		independentExpressions = append(independentExpressions, models.Expression{
			ID:  i,
			Exp: exp,
		})
	}

	return independentExpressions
}
