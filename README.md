# ☁️ Go AWS Lambda Serverless API

Aplicação **Serverless** em **Go (Golang)** desenhada para o runtime `provided.al2023` (ARM64 Graviton2) da **AWS Lambda**, integrada ao **Amazon API Gateway** e acompanhada de definição de **Infraestrutura como Código (IaC)** com AWS SAM (`template.yaml`).

## 🚀 Destaques de Arquitetura Cloud

- **SDK Oficial AWS Lambda Go:** Manipulação tipada de eventos `APIGatewayProxyRequest` e `APIGatewayProxyResponse` com baixíssimo tempo de *cold start*.
- **Infraestrutura como Código (IaC):** Arquivo `template.yaml` pronto para deploy automatizado com **AWS SAM** na arquitetura `arm64`.
- **Simulador Local Integrado:** Detecta automaticamente se está em execução na nuvem (`AWS_LAMBDA_FUNCTION_NAME`) ou localmente, permitindo testar eventos do API Gateway direto no terminal sem custos de nuvem.
- **Testabilidade Total:** Testes unitários validando roteamento de eventos GET/POST e códigos de status HTTP.

## 🛠️ Como Testar e Executar Localmente

    go test -v -race ./...
    go run ./cmd/lambda

## 📦 Como Compilar para Produção na AWS (Linux ARM64)

    GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap ./cmd/lambda
