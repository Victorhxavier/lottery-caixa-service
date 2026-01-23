# 🚀 Como Clonar e Executar

## 📥 Opção 1: Do GitHub (Recomendado)

```bash
# Clonar o repositório
git clone https://github.com/seu-username/lottery-caixa-service.git
cd lottery-caixa-service

# Instalar dependências
go mod download

# Executar
make run-dev
```

## 💻 Opção 2: Local (Arquivos Inclusos)

Se você já tem os arquivos baixados:

```bash
# Navegue até o diretório
cd lottery-caixa-service

# Instale as dependências
go mod download

# Execute
go run ./cmd/main.go
```

## 🐳 Opção 3: Docker Compose (Recomendado para Desenvolvimento)

```bash
# Na raiz do projeto
docker-compose up -d

# Ver logs
docker-compose logs -f lottery-caixa-service

# Parar
docker-compose down
```

## ✅ Verificar Instalação

```bash
# Testar se está funcionando
curl http://localhost:8080/health

# Deve retornar:
# {"status":"healthy","timestamp":"2024-01-15T..."}
```

## 📋 Pré-requisitos

- **Go 1.21+**: https://golang.org/doc/install
- **Git**: https://git-scm.com/download
- **Docker** (opcional): https://docker.com/products/docker-desktop
- **Make** (opcional): `brew install make` (macOS) ou incluso (Linux/Windows)

## 📚 Próximos Passos

1. **Ler**: [QUICKSTART.md](QUICKSTART.md) - 5 minutos para começar
2. **Configurar**: Editar `.env` com suas variáveis
3. **Executar**: `make run-dev` ou `docker-compose up`
4. **Testar**: `curl http://localhost:8080/health`
5. **Explorar**: Ver [README.md](README.md) para documentação completa
6. **Contribuir**: Ver [CONTRIBUTING.md](CONTRIBUTING.md)

## 🔧 Verificar Versão do Go

```bash
go version
# go version go1.21.x ...
```

Se não tiver Go 1.21+, instale de https://golang.org/doc/install

## 🎯 Arquitetura Rápida

```
User/Client
    ↓
[Cloud Run Service] ← você está aqui
    ↓
[Caixa API] (buscar dados)
    ↓
[Downstream Service] (processar dados)
```

## 🚀 Estrutura do Projeto

```
lottery-caixa-service/
├── cmd/main.go                    # Aplicação principal
├── internal/
│   ├── service/service.go         # Lógica de negócio
│   ├── cache/memory.go            # Cache em memória
│   ├── ratelimit/limiter.go       # Rate limiting
│   └── http/handlers/handlers.go  # HTTP endpoints
├── config/config.go               # Configurações
├── Dockerfile                     # Para produção
├── docker-compose.yml             # Desenvolvimento
├── Makefile                       # Automação
└── README.md                      # Doc completa
```

## 🤖 Primeiros Comandos

```bash
# 1. Ver todos os comandos disponíveis
make help

# 2. Instalar dependências
make deps

# 3. Rodar testes
make test

# 4. Executar localmente
make run-dev

# 5. Build do Docker
make docker

# 6. Iniciar stack completa (docker-compose)
make dev
```

## 🔗 URLs Importantes (Desenvolvimento)

| Serviço | URL | Descrição |
|---------|-----|-----------|
| API | http://localhost:8080 | Serviço principal |
| Health | http://localhost:8080/health | Health check |
| Processor | http://localhost:8081 | Serviço downstream |
| Prometheus | http://localhost:9090 | Métricas |
| Grafana | http://localhost:3000 | Visualização (admin/admin) |

## 🌍 Deploy em Produção

Veja [DEPLOY.md](DEPLOY.md) para instruções de:
- Google Cloud Run
- Docker Hub
- Kubernetes
- CI/CD com GitHub Actions

## ⚡ Quick Reference

```bash
# Clonar
git clone <repo> && cd lottery-caixa-service

# Preparar
go mod download

# Desenvolver
make run-dev

# Testar
make test

# Ou Docker
docker-compose up -d
```

## 🐛 Problemas Comuns

### "go: no matching versions"
```bash
go mod tidy
go mod download
```

### "Port 8080 already in use"
```bash
PORT=9000 make run-dev
# ou
PORT=9000 go run ./cmd/main.go
```

### "Docker permission denied"
```bash
sudo usermod -aG docker $USER
# Log out e log in novamente
```

## 📞 Suporte

- **Issues**: Abra uma issue no GitHub
- **Docs**: Veja [README.md](README.md)
- **Quick Start**: [QUICKSTART.md](QUICKSTART.md)

## 🎓 Estrutura de Arquivos

```
go.mod                    # Dependências
go.sum                    # Checksums
cmd/main.go              # Entrypoint
cmd/main_test.go         # Testes
config/config.go         # Variáveis de ambiente
internal/
  service/service.go     # Lógica principal
  cache/memory.go        # Cache
  ratelimit/limiter.go   # Rate limiting
  domain/lottery.go      # Modelos
  http/
    handlers/handlers.go # HTTP handlers
    middleware/middleware.go # Middlewares
examples/processor.go    # Exemplo de processador
scripts/test.sh         # Script de testes
Dockerfile              # Para produção
docker-compose.yml      # Stack local
Makefile                # Automação
README.md               # Documentação
DEPLOY.md               # Deploy guide
QUICKSTART.md           # Quick start
CONTRIBUTING.md         # Guidelines
```

## ✨ Dicas Profissionais

1. **Use Make**: Todos os comandos em `make help`
2. **Docker Compose**: Melhor para desenvolvimento com tudo junto
3. **Logs**: `docker-compose logs -f` para ver logs em tempo real
4. **Testes**: `make test-coverage` para verificar cobertura
5. **Formato**: `make fmt` antes de commits

## 🎯 Fluxo Recomendado

```
1. git clone
   ↓
2. go mod download
   ↓
3. cp .env.example .env (editar se necessário)
   ↓
4. docker-compose up -d
   ↓
5. curl http://localhost:8080/health
   ↓
6. make test (para verificar)
   ↓
7. Começar a desenvolver!
```

---

**Tudo pronto?** Execute `make dev` e comece! 🚀
