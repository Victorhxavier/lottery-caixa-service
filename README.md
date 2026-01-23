# Lottery Caixa Service

Microserviço em Go que consulta resultados da Loteria da Caixa Econômica Federal via Cloud Run e encaminha os dados para um microserviço downstream para processamento.

## 📋 Características

- ✅ Integração com API da Caixa Econômica Federal
- ✅ Rate Limiting (2 req/seg)
- ✅ Caching em memória (5 min TTL)
- ✅ Retry com backoff exponencial
- ✅ Timeout handling robusto
- ✅ Forwarding assincronado para downstream
- ✅ Health checks e readiness probes
- ✅ Métricas Prometheus
- ✅ Estrutura de projeto profissional
- ✅ Testes unitários
- ✅ Docker & Docker Compose
- ✅ Cloud Run ready

## 🚀 Quick Start

### Pré-requisitos

- Go 1.21+
- Docker & Docker Compose (opcional)
- Make (opcional)

### Desenvolvimento Local

```bash
# 1. Clonar repositório
git clone <seu-repositorio>
cd lottery-caixa-service

# 2. Instalar dependências
go mod download

# 3. Copiar arquivo de configuração
cp .env.example .env

# 4. Executar
go run ./cmd/main.go
```

### Com Make

```bash
make help                # Ver todos os comandos
make deps               # Instalar dependências
make build              # Compilar binário
make run-dev            # Executar em desenvolvimento
make test               # Rodar testes
make test-coverage      # Gerar relatório de coverage
make docker             # Build Docker image
make dev                # Iniciar ambiente com docker-compose
```

### Com Docker Compose

```bash
# Iniciar todos os serviços
docker-compose up -d

# Ver logs
docker-compose logs -f lottery-caixa-service

# Parar serviços
docker-compose down
```

## 📁 Estrutura do Projeto

```
lottery-caixa-service/
├── cmd/
│   ├── main.go              # Ponto de entrada
│   └── main_test.go         # Testes
├── config/
│   └── config.go            # Configurações da aplicação
├── internal/
│   ├── domain/
│   │   └── lottery.go       # Modelos de domínio
│   ├── service/
│   │   └── service.go       # Lógica de negócio
│   ├── cache/
│   │   └── memory.go        # Cache em memória
│   ├── ratelimit/
│   │   └── limiter.go       # Rate limiters
│   └── http/
│       ├── handlers/
│       │   └── handlers.go  # HTTP handlers
│       └── middleware/
│           └── middleware.go # Middlewares HTTP
├── go.mod                   # Dependências
├── go.sum                   # Checksums
├── Dockerfile               # Build para produção
├── docker-compose.yml       # Orquestração local
├── Makefile                 # Automação
├── .env.example             # Variáveis de exemplo
├── README.md                # Este arquivo
└── DEPLOY.md               # Guia de deploy
```

## 🔧 Configuração

Todas as configurações podem ser definidas via variáveis de ambiente:

```bash
# Aplicação
APP_NAME=lottery-caixa-service
APP_VERSION=1.0.0
ENVIRONMENT=development
PORT=8080

# API da Caixa
CAIXA_BASE_URL=http://servicebus2.caixa.gov.br/portaldeloterias/api/
CAIXA_TIMEOUT=15
CAIXA_RATE_LIMIT=2
CAIXA_CACHE_TTL=5

# Serviço Downstream
DOWNSTREAM_SERVICE_URL=http://localhost:8081/api/lottery/process
DOWNSTREAM_TIMEOUT=10
ASYNC_FORWARDING=true
```

Veja `.env.example` para todas as opções.

## 📊 Endpoints

### GET /

Busca resultado da loteria

```bash
curl "http://localhost:8080/?gameType=lotofacil"
```

**Parâmetros:**

- `gameType` (opcional): lotofacil, lotomania, duplasena, etc. Default: lotofacil

**Response:**

```json
{
  "status": "success",
  "results": [
    {
      "id": "lotofacil-2024001",
      "gameType": "lotofacil",
      "drawNumber": 2024001,
      "drawDate": "2024-01-15",
      "numbers": [1, 5, 8, 12, 15, 18, 22, 25],
      "winners": 150,
      "prize": 5000.0,
      "processedAt": "2024-01-15T10:30:00Z"
    }
  ],
  "metadata": {
    "processedAt": "2024-01-15T10:30:00Z",
    "sourceService": "lottery-caixa-service",
    "totalRecords": 1,
    "requestId": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

### GET /health

Health check

```bash
curl http://localhost:8080/health
```

### GET /ready

Readiness probe

```bash
curl http://localhost:8080/ready
```

### GET /metrics

Métricas Prometheus

```bash
curl http://localhost:8080/metrics
```

### POST /api/v1/lottery/webhook

Processa webhook de loteria

```bash
curl -X POST http://localhost:8080/api/v1/lottery/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "gameType": "lotofacil",
    "drawNumber": 2024001,
    "drawDate": "2024-01-15",
    "numbers": [1, 5, 8, 12, 15, 18, 22, 25],
    "timestamp": "2024-01-15T10:30:00Z"
  }'
```

## 🧪 Testes

```bash
# Rodar testes
go test -v ./...

# Testes com race detector
go test -v -race ./...

# Com cobertura
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🐳 Docker

### Build

```bash
docker build -t lottery-caixa-service:latest .
```

### Run

```bash
docker run -it --rm \
  -p 8080:8080 \
  -e DOWNSTREAM_SERVICE_URL=http://host.docker.internal:8081 \
  lottery-caixa-service:latest
```

## ☁️ Cloud Run (GCP)

Veja [DEPLOY.md](DEPLOY.md) para instruções completas.

### Quick Deploy

```bash
# Build
gcloud builds submit --tag gcr.io/seu-projeto/lottery-caixa-service

# Deploy
gcloud run deploy lottery-caixa-service \
  --image gcr.io/seu-projeto/lottery-caixa-service \
  --platform managed \
  --region us-central1 \
  --set-env-vars DOWNSTREAM_SERVICE_URL=https://seu-ms/api
```

## 📈 Monitoramento

### Prometheus

- URL: http://localhost:9090
- Métricas: http://localhost:8080/metrics

### Grafana

- URL: http://localhost:3000
- User: admin / Password: admin

## 🔐 Segurança

- Non-root user (UID 1000)
- HTTPS em produção
- Validação de entrada
- Rate limiting
- Timeouts configuráveis
- Health checks

## 🛠️ Troubleshooting

### Erro: "Connection refused"

- Verificar se o downstream está rodando
- Verificar URL em DOWNSTREAM_SERVICE_URL

### Alto latência

- Aumentar CPU em Cloud Run
- Verificar hit rate do cache
- Verificar rate limiting

### Logs

```bash
# Ver logs locais
docker-compose logs -f lottery-caixa-service

# Cloud Logging (GCP)
gcloud logs read --service lottery-caixa-service
```

## 📚 Recursos Adicionais

- [Go Docs](https://golang.org/doc)
- [Cloud Run Docs](https://cloud.google.com/run/docs)
- [Caixa API](https://servicebus.caixa.gov.br)
- [Prometheus Docs](https://prometheus.io)

## 🔄 CI/CD

Exemplos de workflows para GitHub Actions estão em `.github/workflows/`.

## 📝 Logs

O serviço usa logs estruturados:

```
2024-01-15T10:30:00Z [550e8400-e29b-41d4-a716-446655440000] INFO: Buscando resultados para: lotofacil
2024-01-15T10:30:01Z [550e8400-e29b-41d4-a716-446655440000] INFO: Cache hit para lotofacil
2024-01-15T10:30:02Z [550e8400-e29b-41d4-a716-446655440000] INFO: Resultado enviado para downstream
```

## 🤝 Contribuindo

1. Fork o projeto
2. Crie uma branch (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add AmazingFeature'`)
4. Push (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📄 Licença

MIT License - veja LICENSE para detalhes

## ✉️ Suporte

- Issues: [GitHub Issues](https://github.com/seu-repo/issues)
- Email: seu-email@example.com

---

**Desenvolvido com ❤️ em Go**
