package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nayarapoliny/go-aws-lambda-serverless/internal/handler"
)

func main() {
	store := handler.NewMemoryCloudStore()
	router := handler.NewLambdaRouter(store)

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		lambda.Start(router.HandleEvent)
		return
	}

	fmt.Println("☁️ Modo Simulador Local AWS Lambda + API Gateway")
	resp, _ := router.HandleEvent(context.Background(), events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Body:       `{"name":"CloudMetricEvent","category":"Telemetry","payload":"CPU 42%"}`,
	})
	fmt.Printf("➡️ POST /items [Status %d]: %s\n", resp.StatusCode, resp.Body)

	listResp, _ := router.HandleEvent(context.Background(), events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
	})
	fmt.Printf("➡️ GET  /items [Status %d]: %s\n", listResp.StatusCode, listResp.Body)
}
