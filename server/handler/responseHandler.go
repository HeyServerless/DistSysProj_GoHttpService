package handler

import (
	// "encoding/json"

	"net/http"

	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	// appsv1 "k8s.io/api/apps/v1"
	// corev1 "k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/apimachinery/pkg/util/intstr"
	// "k8s.io/client-go/util/retry"
	// "k8s.io/apimachinery/pkg/runtime/schema"
	// "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	models "github.com/goService/server/model"
	services "github.com/goService/server/services"
)

func ResponseFromQueueHandler(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"result": "token",
	})
}

func GetExpressionResultFromDatabase(c *gin.Context) {
	ExpressionResultRequest := models.ExpressionResultRequest{}

	if err := c.ShouldBindJSON(&ExpressionResultRequest); err != nil {
		log.Println("Error: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ReponseObject, _ := services.GetResponseFromDb(c, &models.GetResponseModel{RequestId: ExpressionResultRequest.RequestId})
	if ReponseObject.Response == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"result": "Result is not ready yet",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"result": ReponseObject,
	})

}

func UpdateResultOfExpressionIntoDB(c *gin.Context) {
	fmt.Println("reponse updating")
	ExpressionResponse := models.UpdateResponseRequest{}
	if err := c.ShouldBindJSON(&ExpressionResponse); err != nil {
		log.Println("Error: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	log.Println("ExpressionResponse: ", ExpressionResponse)

	// update the result of the expression into the database
	requestid, err := services.UpdateRequestDataIntoDb(c, &models.UpdateRequestRecordModel{RequestId: ExpressionResponse.RequestId, Response: ExpressionResponse.Response})
	if err != nil {
		log.Println("Error: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"requestId": requestid,
	})

}
