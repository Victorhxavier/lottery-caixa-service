# Quick Start Guide

Comece em 5 minutos!

## ⚡ Iniciação Rápida (Local)

### 1. Clonar Repositório

```bash
git clone https://github.com/seu-username/lottery-caixa-service.git
cd lottery-caixa-service
```

### 2. Instalar Dependências

```bash
go mod download
```

### 3. Configurar Variáveis de Ambiente

```bash
cp .env.example .env
# Editar .env conforme necessário
```

### 4. Executar Localmente

```bash
# Opção 1: Com Go direto
go run ./cmd/main.go

# Opção 2: Com Make
make run-dev

# Opção 3: Build primeiro
make build
./lottery-service
```

### 5. Testar

```bash
curl http://localhost:8080/health
```

## 🐳 Com Docker Compose

```bash
# Iniciar todos os serviços
docker-compose up -d

# Ver logs
docker-compose logs -f lottery-caixa-service

# Parar serviços
docker-compose down
```

## 📊 Acessar Serviços

- **API**: http://localhost:8080
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin)

## 🎯 Endpoints Principais

```bash
# Health check
curl http://localhost:8080/health

# Buscar resultado de loteria
curl "http://localhost:8080/?gameType=lotofacil"

# Enviar webhook
curl -X POST http://localhost:8080/api/v1/lottery/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "gameType": "lotofacil",
    "drawNumber": 2024001,
    "drawDate": "2024-01-15",
    "numbers": [1,5,8,12,15,18,22,25]
  }'
```

## 🧪 Rodar Testes

```bash
# Testes simples
make test

# Com cobertura
make test-coverage

# Específico
go test -v -run TestHealthCheck ./...
```

## 🔨 Make Commands

```bash
make help              # Ver todos os comandos
make deps             # Instalar dependências
make build            # Compilar binário
make run              # Executar (após build)
make run-dev          # Executar com Go direto
make test             # Rodar testes
make test-coverage    # Testes com coverage
make fmt              # Formatar código
make lint             # Rodar linter
make docker           # Build Docker image
make dev              # Iniciar docker-compose
make dev-logs         # Ver logs do docker-compose
make dev-down         # Parar docker-compose
make clean            # Limpar artefatos
```

## 📁 Estrutura

```
lottery-caixa-service/
├── cmd/main.go              # Entrada principal
├── internal/
│   ├── service/service.go   # Lógica principal
│   ├── cache/memory.go      # Cache
│   ├── ratelimit/limiter.go # Rate limiting
│   └── http/handlers/       # HTTP handlers
├── config/config.go         # Configuração
├── Dockerfile               # Para produção
├── docker-compose.yml       # Para desenvolvimento
├── Makefile                 # Automação
└── README.md               # Documentação completa
```

## ⚙️ Configuração Essencial

As principais variáveis de ambiente:

```bash
PORT=8080                                          # Porta
ENVIRONMENT=development                            # Ambiente
DOWNSTREAM_SERVICE_URL=http://localhost:8081     # URL downstream
CAIXA_RATE_LIMIT=2                                # Req/seg
CAIXA_CACHE_TTL=5                                 # Cache TTL em minutos
```

Ver `.env.example` para todas as opções.

## 🔍 Logs

```bash
# Stderr (logs estruturados)
# Aparecem no terminal quando executa go run

# Docker Compose
docker-compose logs -f lottery-caixa-service
docker-compose logs -f lottery-processor

# Cloud Run
gcloud logs read --service lottery-caixa-service
```

## 🐛 Troubleshooting Comum

### Erro: "connection refused"
```bash
# Verificar se a porta 8080 está em uso
lsof -i :8080
# Ou especificar outra porta
PORT=9000 go run ./cmd/main.go
```

### Erro: "go: no matching versions"
```bash
# Atualizar go.mod
go get -u ./...
go mod tidy
```

### Docker: "permission denied"
```bash
# Dar permissão
sudo usermod -aG docker $USER
# Ou executar com sudo
sudo docker-compose up
```

## 💡 Próximos Passos

1. **Explore o código**: Veja `internal/service/service.go`
2. **Customize**: Edite configurações em `.env`
3. **Estenda**: Adicione novos handlers em `internal/http/handlers/`
4. **Deploy**: Veja [DEPLOY.md](DEPLOY.md) para Cloud Run

## 📚 Recursos

- [README.md](README.md) - Documentação completa
- [DEPLOY.md](DEPLOY.md) - Deploy em produção
- [CONTRIBUTING.md](CONTRIBUTING.md) - Como contribuir
- [Go Docs](https://golang.org/doc)
- [Cloud Run Docs](https://cloud.google.com/run/docs)

## 🆘 Precisa de Ajuda?

1. Verifique o README.md
2. Veja os exemplos em `examples/`
3. Abra uma issue no GitHub
4. Veja logs: `docker-compose logs -f`

---

**Tudo configurado? Comece a desenvolver!** 🚀

```bash
make dev    # Inicia tudo
```
