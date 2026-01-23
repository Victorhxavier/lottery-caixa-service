# Guia de Execução Local

Este guia fornece instruções passo a passo para executar o Lottery Caixa Service localmente em sua máquina.

## 📋 Pré-requisitos

- **Go 1.21+** (você tem: `go version go1.23.2`)
- **Git** (para clonar o repositório)
- **Docker & Docker Compose** (opcional, para ambiente completo)

## 🚀 Opções de Execução

### Opção 1: Script Automático (Recomendado)

O jeito mais fácil de começar:

```bash
./start-local.sh
```

O script oferece 3 modos:
1. **Apenas serviço principal** - Porta 8080
2. **Serviço + Processador** - Portas 8080 e 8081
3. **Ambiente completo** - Docker Compose com todos os serviços

### Opção 2: Execução Manual Simples

```bash
# 1. Baixar dependências (se ainda não fez)
go mod download

# 2. Executar o serviço
go run ./cmd/main.go
```

O serviço estará disponível em `http://localhost:8080`

### Opção 3: Com Processador Downstream

**Terminal 1 - Processador:**
```bash
cd /Users/victorhxavier/Documents/lottery-caixa-service
PORT=8081 go run ./examples/processor.go
```

**Terminal 2 - Serviço Principal:**
```bash
cd /Users/victorhxavier/Documents/lottery-caixa-service
go run ./cmd/main.go
```

### Opção 4: Com Make

```bash
# Ver comandos disponíveis
make help

# Executar em modo dev
make run-dev

# Ou compilar e executar
make build
./lottery-service
```

### Opção 5: Docker Compose (Ambiente Completo)

```bash
# Iniciar todos os serviços
docker-compose up

# Ou em background
docker-compose up -d

# Ver logs
docker-compose logs -f lottery-caixa-service

# Parar
docker-compose down
```

## 🔧 Configuração

O arquivo `.env` já foi criado com as configurações padrão. Você pode editá-lo conforme necessário:

```bash
# Editar configurações
nano .env
```

### Configurações Principais

```bash
PORT=8080                                      # Porta do serviço
ENVIRONMENT=development                        # Ambiente
DOWNSTREAM_SERVICE_URL=http://localhost:8081  # URL do processador
CAIXA_RATE_LIMIT=2                            # Requisições por segundo
CAIXA_CACHE_TTL=5                             # Cache TTL em minutos
```

## 🧪 Testando o Serviço

### 1. Health Check

```bash
curl http://localhost:8080/health
```

Resposta esperada:
```json
{
  "status": "healthy",
  "timestamp": "2026-01-23T..."
}
```

### 2. Buscar Resultado de Loteria

```bash
curl "http://localhost:8080/?gameType=lotofacil"
```

### 3. Enviar Webhook

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

### 4. Ver Métricas

```bash
curl http://localhost:8080/metrics
```

## 📊 Acessando Serviços

Quando executando com Docker Compose:

| Serviço | URL | Credenciais |
|---------|-----|-------------|
| API Principal | http://localhost:8080 | - |
| Processador | http://localhost:8081 | - |
| Prometheus | http://localhost:9090 | - |
| Grafana | http://localhost:3000 | admin/admin |

## 🐛 Troubleshooting

### Erro: "port 8080: bind: address already in use"

Algum processo já está usando a porta 8080.

**Solução 1:** Encontrar e parar o processo
```bash
# macOS/Linux
lsof -ti:8080 | xargs kill -9

# Ou identificar o processo
lsof -i :8080
```

**Solução 2:** Usar outra porta
```bash
PORT=9000 go run ./cmd/main.go
```

### Erro: "connection refused" ao chamar downstream

O processador downstream não está rodando.

**Solução:** Iniciar o processador em outro terminal
```bash
PORT=8081 go run ./examples/processor.go
```

Ou desabilitar o forwarding no `.env`:
```bash
ASYNC_FORWARDING=false
```

### Erro: "checksum mismatch" ao baixar dependências

O arquivo `go.sum` está corrompido.

**Solução:**
```bash
rm go.sum
go mod download
go mod tidy
```

### Go não encontrado

**Solução:** Instalar Go
```bash
# macOS (com Homebrew)
brew install go

# Ou baixe de: https://go.dev/dl/
```

## 🧰 Comandos Úteis

```bash
# Rodar testes
go test -v ./...

# Testes com coverage
make test-coverage

# Formatar código
go fmt ./...

# Ver logs estruturados
go run ./cmd/main.go 2>&1 | grep "INFO\|ERROR"

# Compilar binário
go build -o lottery-service ./cmd/main.go

# Limpar cache
go clean -cache -testcache

# Verificar dependências
go mod verify
```

## 📝 Arquivos Importantes

- `.env` - Configurações locais (criado automaticamente)
- `go.mod` - Dependências do projeto
- `cmd/main.go` - Ponto de entrada do serviço
- `examples/processor.go` - Processador downstream de exemplo
- `docker-compose.yml` - Configuração do ambiente completo

## 🔄 Fluxo de Desenvolvimento

1. **Fazer mudanças no código**
2. **Salvar arquivos**
3. **Reiniciar serviço** (Ctrl+C e rodar novamente)
4. **Testar** com curl ou Postman

### Hot Reload (Opcional)

Para reinicialização automática ao salvar arquivos:

```bash
# Instalar air
go install github.com/cosmtrek/air@latest

# Executar com hot reload
air
```

## 🎯 Próximos Passos

- ✅ Serviço rodando localmente
- 📖 Ler [README.md](README.md) para documentação completa
- 🚀 Ver [DEPLOY.md](DEPLOY.md) para deploy em produção
- 🤝 Ver [CONTRIBUTING.md](CONTRIBUTING.md) para contribuir

## 💡 Dicas

1. **Logs Estruturados**: Os logs aparecem no stderr com timestamps e request IDs
2. **Cache**: O cache padrão é de 5 minutos para economizar chamadas à API
3. **Rate Limiting**: Por padrão, máximo 2 requisições por segundo à API da Caixa
4. **Async Forwarding**: Os dados são enviados ao downstream de forma assíncrona

## 📞 Suporte

- **Documentação**: [README.md](README.md)
- **Quick Start**: [QUICKSTART.md](QUICKSTART.md)
- **Issues**: GitHub Issues

---

**Pronto para desenvolver! 🎉**

Se tiver algum problema, verifique os logs e a seção de troubleshooting acima.
