# Deploy no Google Cloud Run

Guia completo para fazer deploy do Lottery Caixa Service no Google Cloud Run.

## 📋 Pré-requisitos

### 1. Conta GCP
- Conta Google Cloud ativa
- Projeto GCP criado
- Billing habilitado no projeto

### 2. Google Cloud SDK
```bash
# macOS (Homebrew)
brew install --cask google-cloud-sdk

# Ou download direto
# https://cloud.google.com/sdk/docs/install
```

### 3. Autenticação
```bash
# Login na GCP
gcloud auth login

# Configurar projeto padrão
gcloud config set project SEU_PROJECT_ID

# Listar projetos disponíveis
gcloud projects list
```

## 🚀 Deploy Automático (Recomendado)

### Opção 1: Script Interativo

```bash
./deploy-cloudrun.sh
```

O script irá:
1. ✅ Verificar instalação do gcloud
2. ✅ Solicitar Project ID (se não definido)
3. ✅ Habilitar APIs necessárias
4. ✅ Fazer build da imagem Docker
5. ✅ Fazer deploy no Cloud Run
6. ✅ Retornar URL do serviço

### Opção 2: Com Variáveis de Ambiente

```bash
# Definir variáveis
export GCP_PROJECT_ID="seu-projeto-id"
export SERVICE_NAME="lottery-caixa-service"
export REGION="us-central1"
export MEMORY="512Mi"
export CPU="1"
export MIN_INSTANCES="0"
export MAX_INSTANCES="10"

# Executar deploy
./deploy-cloudrun.sh
```

## 🔧 Deploy Manual

### 1. Habilitar APIs

```bash
gcloud services enable cloudbuild.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable containerregistry.googleapis.com
```

### 2. Build da Imagem

**Opção A: Cloud Build (Recomendado)**
```bash
gcloud builds submit --tag gcr.io/SEU_PROJECT_ID/lottery-caixa-service
```

**Opção B: Docker Local**
```bash
# Build local
docker build -t gcr.io/SEU_PROJECT_ID/lottery-caixa-service .

# Configurar Docker para GCR
gcloud auth configure-docker

# Push para GCR
docker push gcr.io/SEU_PROJECT_ID/lottery-caixa-service
```

### 3. Deploy no Cloud Run

**Comando Básico:**
```bash
gcloud run deploy lottery-caixa-service \
  --image gcr.io/SEU_PROJECT_ID/lottery-caixa-service \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated
```

**Comando Completo (com todas as configurações):**
```bash
gcloud run deploy lottery-caixa-service \
  --image gcr.io/SEU_PROJECT_ID/lottery-caixa-service \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --memory 512Mi \
  --cpu 1 \
  --timeout 60s \
  --min-instances 0 \
  --max-instances 10 \
  --set-env-vars "ENVIRONMENT=production,PORT=8080,CAIXA_BASE_URL=https://servicebus2.caixa.gov.br/portaldeloterias/api,CAIXA_TIMEOUT=15,CAIXA_RATE_LIMIT=2,CAIXA_CACHE_TTL=5,CAIXA_MAX_RETRIES=3,CAIXA_RETRY_BACKOFF=1000,ASYNC_FORWARDING=false"
```

### 4. Usando Arquivo YAML

```bash
# Editar cloudrun.yaml e substituir PROJECT_ID
sed -i '' 's/PROJECT_ID/seu-projeto-id/g' cloudrun.yaml

# Deploy usando YAML
gcloud run services replace cloudrun.yaml --region us-central1
```

## 🔐 Configurações de Segurança

### Autenticação (Requer Token)

```bash
# Deploy com autenticação
gcloud run deploy lottery-caixa-service \
  --image gcr.io/SEU_PROJECT_ID/lottery-caixa-service \
  --platform managed \
  --region us-central1 \
  --no-allow-unauthenticated

# Testar com token
curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
  https://seu-servico.run.app/health
```

### Service Account Customizada

```bash
# Criar service account
gcloud iam service-accounts create lottery-service-sa \
  --display-name "Lottery Service Account"

# Deploy com service account
gcloud run deploy lottery-caixa-service \
  --image gcr.io/SEU_PROJECT_ID/lottery-caixa-service \
  --service-account lottery-service-sa@SEU_PROJECT_ID.iam.gserviceaccount.com \
  --region us-central1
```

## 📊 Configurações de Recursos

### CPU e Memória

```bash
# Mínimo (economia)
--memory 256Mi --cpu 1

# Recomendado (produção)
--memory 512Mi --cpu 1

# Alto desempenho
--memory 1Gi --cpu 2
```

### Autoscaling

```bash
# Sempre com instância ativa (sem cold start)
--min-instances 1 --max-instances 10

# Scale to zero (economia máxima)
--min-instances 0 --max-instances 100

# Fixo (sem autoscaling)
--min-instances 5 --max-instances 5
```

### Concorrência

```bash
# Máximo de requisições por container
--concurrency 80

# Sem limite (padrão)
--concurrency 1000
```

## 🌍 Variáveis de Ambiente

### Definir Variáveis

```bash
# Via comando
gcloud run services update lottery-caixa-service \
  --set-env-vars KEY1=value1,KEY2=value2

# Usando arquivo
gcloud run services update lottery-caixa-service \
  --env-vars-file .env.cloudrun
```

### Variáveis Secretas (Secrets Manager)

```bash
# Criar secret
echo -n "valor-secreto" | gcloud secrets create my-secret --data-file=-

# Usar no Cloud Run
gcloud run services update lottery-caixa-service \
  --update-secrets DATABASE_PASSWORD=my-secret:latest
```

## 🔍 Monitoramento

### Ver Logs

```bash
# Logs em tempo real
gcloud run logs tail lottery-caixa-service --region us-central1

# Logs das últimas 2 horas
gcloud run logs read lottery-caixa-service \
  --region us-central1 \
  --limit 100

# Filtrar por severidade
gcloud run logs read lottery-caixa-service \
  --region us-central1 \
  --log-filter 'severity>=ERROR'
```

### Métricas

```bash
# Ver estatísticas
gcloud run services describe lottery-caixa-service \
  --region us-central1 \
  --format yaml

# Console de métricas
# https://console.cloud.google.com/run/detail/REGION/SERVICE_NAME/metrics
```

## 🧪 Testando o Serviço

```bash
# Obter URL
SERVICE_URL=$(gcloud run services describe lottery-caixa-service \
  --platform managed \
  --region us-central1 \
  --format 'value(status.url)')

# Health check
curl $SERVICE_URL/health

# Buscar loteria
curl "$SERVICE_URL/api/v1/lottery?gameType=megasena"

# Métricas
curl $SERVICE_URL/metrics
```

## 🔄 Atualizações

### Deploy de Nova Versão

```bash
# 1. Build nova imagem
gcloud builds submit --tag gcr.io/SEU_PROJECT_ID/lottery-caixa-service:v2

# 2. Deploy com tag específica
gcloud run deploy lottery-caixa-service \
  --image gcr.io/SEU_PROJECT_ID/lottery-caixa-service:v2 \
  --region us-central1

# Ou usar o script
./deploy-cloudrun.sh
```

### Rollback

```bash
# Listar revisões
gcloud run revisions list --service lottery-caixa-service --region us-central1

# Rotear tráfego para revisão anterior
gcloud run services update-traffic lottery-caixa-service \
  --to-revisions REVISION_NAME=100 \
  --region us-central1
```

### Traffic Splitting (Canary/Blue-Green)

```bash
# 50% em cada revisão
gcloud run services update-traffic lottery-caixa-service \
  --to-revisions REVISION_NEW=50,REVISION_OLD=50 \
  --region us-central1

# Migrar gradualmente
gcloud run services update-traffic lottery-caixa-service \
  --to-revisions REVISION_NEW=10,REVISION_OLD=90 \
  --region us-central1
```

## 💰 Custos Estimados

### Calculadora
- Free tier: 2 milhões de requisições/mês
- 360,000 GB-seconds de memória
- 180,000 vCPU-seconds

### Exemplo de Uso
**Cenário:** 1M requisições/mês, 200ms latência média, 512MB RAM

```
- Requisições: GRÁTIS (dentro do free tier)
- CPU: ~$2/mês
- Memória: ~$1/mês
- Total: ~$3/mês
```

## 🛠️ Troubleshooting

### Erro: "Permission denied"

```bash
# Dar permissões necessárias
gcloud projects add-iam-policy-binding SEU_PROJECT_ID \
  --member="user:seu-email@gmail.com" \
  --role="roles/run.admin"
```

### Erro: "Quota exceeded"

```bash
# Aumentar quota no console
# https://console.cloud.google.com/iam-admin/quotas
```

### Erro: "Container failed to start"

```bash
# Ver logs detalhados
gcloud run logs read lottery-caixa-service --region us-central1 --limit 50

# Verificar health check
curl https://seu-servico.run.app/health
```

### Cold Start Alto

```bash
# Configurar min-instances
gcloud run services update lottery-caixa-service \
  --min-instances 1 \
  --region us-central1
```

## 📚 Recursos Úteis

- [Cloud Run Docs](https://cloud.google.com/run/docs)
- [Pricing Calculator](https://cloud.google.com/products/calculator)
- [Best Practices](https://cloud.google.com/run/docs/tips)
- [Console Cloud Run](https://console.cloud.google.com/run)

## 🔗 Comandos Úteis

```bash
# Listar serviços
gcloud run services list

# Deletar serviço
gcloud run services delete lottery-caixa-service --region us-central1

# Ver configuração
gcloud run services describe lottery-caixa-service --region us-central1

# Pausar serviço (min-instances=0)
gcloud run services update lottery-caixa-service --min-instances 0

# Ver custos
gcloud billing accounts list
```

## 🎯 Checklist de Deploy

- [ ] Projeto GCP criado e billing habilitado
- [ ] Google Cloud SDK instalado e autenticado
- [ ] APIs habilitadas (Cloud Run, Cloud Build, Container Registry)
- [ ] Dockerfile testado localmente
- [ ] Variáveis de ambiente configuradas
- [ ] Script de deploy executado com sucesso
- [ ] Health check funcionando
- [ ] Endpoints testados
- [ ] Monitoramento configurado
- [ ] Alertas configurados (opcional)
- [ ] Custom domain configurado (opcional)

---

**Pronto para produção! 🚀**
