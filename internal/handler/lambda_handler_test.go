package handler

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestLambdaRouter_PostAndGet(t *testing.T) {
	store := NewMemoryCloudStore()
	router := NewLambdaRouter(store)
	ctx := context.Background()

	postReq := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodPost,
		Body:       `{"name":"OrderWebhook","category":"SQS","payload":"processed"}`,
	}

	postResp, err := router.HandleEvent(ctx, postReq)
	if err != nil || postResp.StatusCode != http.StatusCreated {
		t.Fatalf("esperava status 201 Created, obteve %d (err: %v)", postResp.StatusCode, err)
	}

	getReq := events.APIGatewayProxyRequest{
		HTTPMethod:            http.MethodGet,
		QueryStringParameters: map[string]string{"id": "LAMBDA-ITEM-002"},
	}

	getResp, err := router.HandleEvent(ctx, getReq)
	if err != nil || getResp.StatusCode != http.StatusOK {
		t.Fatalf("esperava encontrar o item recém-criado com status 200, obteve %d", getResp.StatusCode)
	}
}
