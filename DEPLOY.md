f# Deploy Guide - Lottery Caixa Service

Guia completo para fazer deploy do serviço em diferentes ambientes.

## 📋 Pré-requisitos

- Google Cloud Account
- `gcloud` CLI instalado
- Docker instalado
- Git

## 🚀 Deploy no Google Cloud Run

### 1. Setup Inicial

```bash
# Login no Google Cloud
gcloud auth login

# Configurar projeto
gcloud config set project seu-projeto-id

# Habilitar APIs necessárias
gcloud services enable run.googleapis.com
gcloud services enable cloudbuild.googleapis.com
gcloud services enable containerregistry.googleapis.com
```

### 2. Build da Imagem Docker

```bash
# Build local
docker build -t lottery-caixa-service:latest .

# Tag para GCP
docker tag lottery-caixa-service:latest \
  gcr.io/mega-sena-79491/lottery-caixa-service:latest

# Push para Container Registry
docker push gcr.io/mega-sena-79491/lottery-caixa-service:latest
```

### 3. Deploy no Cloud Run

#### Via CLI

```bash
gcloud run deploy lottery-caixa-service \
  --image gcr.io/seu-projeto-id/lottery-caixa-service:latest \
  --platform managed \
  --region us-central1 \
  --memory 512Mi \
  --cpu 1 \
  --timeout 300 \
  --max-instances 10 \
  --set-env-vars \
    ENVIRONMENT=production,\
    CAIXA_RATE_LIMIT=2,\
    CAIXA_CACHE_TTL=5,\
    DOWNSTREAM_SERVICE_URL=https://seu-ms-receptor/api/lottery/process \
  --allow-unauthenticated
```

#### Via Google Cloud Console

1. Vá para Cloud Run
2. Click "Create Service"
3. Selecione a imagem: `gcr.io/seu-projeto-id/lottery-caixa-service:latest`
4. Configure:
   - **Service name**: lottery-caixa-service
   - **Region**: us-central1
   - **Memory**: 512 MB
   - **CPU**: 1
   - **Timeout**: 300s
   - **Max instances**: 10
5. Clique "Create"

### 4. Configurar Variáveis de Ambiente

Na aba "Runtime settings" do serviço:

```
ENVIRONMENT=production
DOWNSTREAM_SERVICE_URL=https://seu-ms-receptor/api/lottery/process
CAIXA_RATE_LIMIT=2
CAIXA_CACHE_TTL=5
```

### 5. Verificar Deploy

```bash
# Obter URL do serviço
gcloud run services describe lottery-caixa-service \
  --platform managed \
  --region us-central1 \
  --format='value(status.url)'

# Testar health check
curl https://seu-cloud-run-url/health

# Testar endpoint
curl "https://seu-cloud-run-url/?gameType=lotofacil"
```

## 📊 Monitoramento

### Cloud Logging

```bash
# Ver logs em tempo real
gcloud logs read \
  "resource.type=cloud_run_revision AND \
   resource.labels.service_name=lottery-caixa-service" \
  --limit 50 \
  --follow

# Logs com filtro de erro
gcloud logs read \
  "resource.type=cloud_run_revision AND \
   resource.labels.service_name=lottery-caixa-service AND \
   severity=ERROR"
```

### Cloud Monitoring

1. Vá para Cloud Monitoring
2. Create a new dashboard
3. Adicione métricas:
   - `cloud_run_request_count`
   - `cloud_run_request_latencies`
   - `cloud_run_instance_cpu_usage`
   - `cloud_run_instance_memory_usage`

## 🔐 Segurança

### Service Account

```bash
# Criar service account
gcloud iam service-accounts create lottery-caixa-service \
  --display-name="Lottery Caixa Service"

# Grant roles necessários
gcloud projects add-iam-policy-binding seu-projeto-id \
  --member=serviceAccount:lottery-caixa-service@seu-projeto-id.iam.gserviceaccount.com \
  --role=roles/logging.logWriter
```

### Secret Manager

Para credenciais sensíveis:

```bash
# Criar secret
echo -n "seu-token" | gcloud secrets create CAIXA_API_TOKEN \
  --data-file=-

# Usar em Cloud Run
gcloud run deploy lottery-caixa-service \
  --update-secrets CAIXA_API_TOKEN=CAIXA_API_TOKEN:latest
```

## 🔄 CI/CD com GitHub Actions

Criar `.github/workflows/deploy.yml`:

```yaml
name: Deploy to Cloud Run

on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v2
    
    - name: Setup Cloud SDK
      uses: google-github-actions/setup-gcloud@v0
      with:
        service_account_key: ${{ secrets.GCP_SA_KEY }}
        project_id: ${{ secrets.GCP_PROJECT_ID }}
    
    - name: Configure Docker for GCR
      run: gcloud auth configure-docker
    
    - name: Build and push Docker image
      run: |
        docker build -t gcr.io/${{ secrets.GCP_PROJECT_ID }}/lottery-caixa-service:${{ github.sha }} .
        docker push gcr.io/${{ secrets.GCP_PROJECT_ID }}/lottery-caixa-service:${{ github.sha }}
    
    - name: Deploy to Cloud Run
      run: |
        gcloud run deploy lottery-caixa-service \
          --image gcr.io/${{ secrets.GCP_PROJECT_ID }}/lottery-caixa-service:${{ github.sha }} \
          --platform managed \
          --region us-central1 \
          --allow-unauthenticated
```

## 📈 Scaling

### Configurar Auto Scaling

```bash
gcloud run update lottery-caixa-service \
  --max-instances 100 \
  --min-instances 1 \
  --concurrency 100 \
  --region us-central1
```

### Aumentar Recursos

```bash
gcloud run update lottery-caixa-service \
  --memory 1Gi \
  --cpu 2 \
  --region us-central1
```

## 🔗 VPC Connector

Para conectar com serviços internos:

```bash
# Criar VPC Connector
gcloud compute networks vpc-access connectors create lottery-connector \
  --region us-central1 \
  --subnet default

# Deploy com VPC Connector
gcloud run deploy lottery-caixa-service \
  --image gcr.io/seu-projeto-id/lottery-caixa-service:latest \
  --vpc-connector lottery-connector \
  --region us-central1
```

## 📦 Versioning

### Blue-Green Deployment

```bash
# Deploy nova versão como draft
gcloud run deploy lottery-caixa-service-v2 \
  --image gcr.io/seu-projeto-id/lottery-caixa-service:v2 \
  --region us-central1 \
  --no-traffic

# Testar nova versão
curl https://lottery-caixa-service-v2-xxxxx.run.app/health

# Migrar tráfego
gcloud run services update-traffic lottery-caixa-service \
  --to-revisions LATEST=100
```

## 🛑 Rollback

```bash
# Listar revisions
gcloud run revisions list --service lottery-caixa-service

# Reverter para revision anterior
gcloud run services update-traffic lottery-caixa-service \
  --to-revisions lottery-caixa-service-xxxxx=100
```

## 📊 Performance Tuning

### Otimizações

1. **Memory**: Aumentar para 512Mi-1Gi para melhor performance
2. **CPU**: Usar CPU 2 para workloads cpu-bound
3. **Concurrency**: Aumentar para 100-1000
4. **Min Instances**: Usar 1 para evitar cold starts

### Exemplo Otimizado

```bash
gcloud run deploy lottery-caixa-service \
  --image gcr.io/seu-projeto-id/lottery-caixa-service:latest \
  --memory 1Gi \
  --cpu 2 \
  --concurrency 500 \
  --min-instances 2 \
  --max-instances 50 \
  --timeout 600 \
  --region us-central1
```

## 🧪 Testes em Staging

```bash
# Deploy em staging
gcloud run deploy lottery-caixa-service-staging \
  --image gcr.io/seu-projeto-id/lottery-caixa-service:latest \
  --region us-central1

# Testar
curl "https://lottery-caixa-service-staging-xxxxx.run.app/?gameType=lotofacil"

# Deletar após testes
gcloud run services delete lottery-caixa-service-staging --region us-central1
```

## 📝 Troubleshooting

### Erro: "Failed to pull image"

```bash
# Verificar permissões
gcloud projects get-iam-policy seu-projeto-id \
  --flatten="bindings[].members" \
  --filter="bindings.role:roles/container.developer"
```

### Erro: "Service Unavailable"

```bash
# Ver logs de erro
gcloud logs read "resource.type=cloud_run_revision" --limit 100

# Aumentar timeout
gcloud run update lottery-caixa-service --timeout 600
```

### Alto latência

```bash
# Aumentar CPU
gcloud run update lottery-caixa-service --cpu 2

# Aumentar memory
gcloud run update lottery-caixa-service --memory 1Gi

# Aumentar min instances
gcloud run update lottery-caixa-service --min-instances 2
```

## 🎯 Checklist de Deploy

- [ ] Testar localmente
- [ ] Aumentar versão em go.mod
- [ ] Atualizar README.md
- [ ] Build Docker image
- [ ] Push para Container Registry
- [ ] Deploy em staging
- [ ] Testar em staging
- [ ] Deploy em produção
- [ ] Verificar logs
- [ ] Monitorar métricas

## 📚 Recursos Adicionais

- [Cloud Run Best Practices](https://cloud.google.com/run/docs/quickstarts/build-and-deploy)
- [Cloud Run Security](https://cloud.google.com/run/docs/securing)
- [Cloud Run Limits](https://cloud.google.com/run/quotas)

---

**Última atualização**: 2024-01-15
