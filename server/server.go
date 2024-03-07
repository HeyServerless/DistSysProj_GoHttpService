package server

import (
	"fmt"

	"github.com/gin-gonic/gin"

	// "honnef.co/go/tools/conf

	"flag"

	"github.com/goService/server/handler"
	services "github.com/goService/server/services"

	"k8s.io/client-go/kubernetes"
	// pb "github.com/example/mypackage"
	// EKS requires the AWS SDK to be imported
)

// Server is srv struct that holds srv Kubernetes client
type Server struct {
	KubeClient *kubernetes.Clientset
}

func ApiMiddleware(cli *kubernetes.Clientset) gin.HandlerFunc {
	// do something with the request
	return func(c *gin.Context) {
		// do something with the request

		c.Set("kubeClient", cli)
		c.Next()
	}
}

func (srv *Server) Initialize() {

	flag.Parse()

	CreateInboundQueue := services.CreateQueue("InboundQueue", "https://5ycsge77e1.execute-api.us-east-1.amazonaws.com/default/sqsCallBackFunction")
	CreateOutbounQueue := services.CreateQueue("OutboundQueue", "https://5ycsge77e1.execute-api.us-east-1.amazonaws.com/default/sqsCallBackFunction")
	fmt.Println(CreateInboundQueue)
	fmt.Println(CreateOutbounQueue)

	fmt.Println("=================================starting server=================================")
	r := gin.Default()
	// r.Use(ApiMiddleware(client))

	/** temp routes end*/

	// main service routes
	r.POST("/calculateExpression", handler.HandleInputExpression)
	r.POST("/updateResultOfExpression", handler.UpdateResultOfExpressionIntoDB)

	r.POST("/getExpressionResult", handler.GetExpressionResultFromDatabase)

	r.Run(":8080")
}
