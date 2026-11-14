package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

type CloudItem struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	Payload   string    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

type Storage interface {
	Save(ctx context.Context, item *CloudItem) error
	GetByID(ctx context.Context, id string) (*CloudItem, bool)
	ListAll(ctx context.Context) []CloudItem
}

type MemoryCloudStore struct {
	mu      sync.RWMutex
	items   map[string]CloudItem
	counter uint64
}

func NewMemoryCloudStore() *MemoryCloudStore {
	store := &MemoryCloudStore{
		items: make(map[string]CloudItem),
	}
	_ = store.Save(context.Background(), &CloudItem{
		Name:     "Serverless Event Processor",
		Category: "AWS-Lambda",
		Payload:  "Cold-start optimized Go runtime",
	})
	return store
}

func (s *MemoryCloudStore) Save(_ context.Context, item *CloudItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := atomic.AddUint64(&s.counter, 1)
	item.ID = fmt.Sprintf("LAMBDA-ITEM-%03d", id)
	item.CreatedAt = time.Now().UTC()
	s.items[item.ID] = *item
	return nil
}

func (s *MemoryCloudStore) GetByID(_ context.Context, id string) (*CloudItem, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.items[id]
	return &val, ok
}

func (s *MemoryCloudStore) ListAll(_ context.Context) []CloudItem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]CloudItem, 0, len(s.items))
	for _, v := range s.items {
		list = append(list, v)
	}
	return list
}

type LambdaRouter struct {
	store Storage
}

func NewLambdaRouter(store Storage) *LambdaRouter {
	return &LambdaRouter{store: store}
}

func (l *LambdaRouter) HandleEvent(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case http.MethodGet:
		if id := req.QueryStringParameters["id"]; id != "" {
			item, ok := l.store.GetByID(ctx, id)
			if !ok {
				return jsonResponse(http.StatusNotFound, map[string]string{"error": "item não encontrado na Lambda"}), nil
			}
			return jsonResponse(http.StatusOK, item), nil
		}
		return jsonResponse(http.StatusOK, l.store.ListAll(ctx)), nil

	case http.MethodPost:
		var item CloudItem
		if err := json.Unmarshal([]byte(req.Body), &item); err != nil {
			return jsonResponse(http.StatusBadRequest, map[string]string{"error": "corpo JSON inválido"}), nil
		}
		if strings.TrimSpace(item.Name) == "" {
			return jsonResponse(http.StatusBadRequest, map[string]string{"error": "o campo 'name' é obrigatório"}), nil
		}
		_ = l.store.Save(ctx, &item)
		return jsonResponse(http.StatusCreated, item), nil

	default:
		return jsonResponse(http.StatusMethodNotAllowed, map[string]string{"error": "método HTTP não suportado"}), nil
	}
}

func jsonResponse(status int, body interface{}) events.APIGatewayProxyResponse {
	raw, _ := json.Marshal(body)
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(raw),
	}
}
