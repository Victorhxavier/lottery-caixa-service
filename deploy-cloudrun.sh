#!/bin/bash

# Script de deploy para Google Cloud Run
# Author: Claude Code

set -e

# Cores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo "========================================="
echo "  Deploy para Google Cloud Run"
echo "========================================="
echo ""

# Configurações (edite conforme necessário)
PROJECT_ID="${GCP_PROJECT_ID:-}"
SERVICE_NAME="${SERVICE_NAME:-lottery-caixa-service}"
REGION="${REGION:-us-central1}"
IMAGE_NAME="gcr.io/${PROJECT_ID}/${SERVICE_NAME}"
MIN_INSTANCES="${MIN_INSTANCES:-0}"
MAX_INSTANCES="${MAX_INSTANCES:-10}"
MEMORY="${MEMORY:-512Mi}"
CPU="${CPU:-1}"
TIMEOUT="${TIMEOUT:-60}"

# Verificar se gcloud está instalado
if ! command -v gcloud &> /dev/null; then
    echo -e "${RED}✗ Google Cloud SDK (gcloud) não está instalado.${NC}"
    echo "Instale em: https://cloud.google.com/sdk/docs/install"
    exit 1
fi

# Verificar se PROJECT_ID está definido
if [ -z "$PROJECT_ID" ]; then
    echo -e "${YELLOW}⚠️  GCP_PROJECT_ID não definido.${NC}"
    read -p "Digite o ID do seu projeto GCP: " PROJECT_ID

    if [ -z "$PROJECT_ID" ]; then
        echo -e "${RED}✗ PROJECT_ID é obrigatório${NC}"
        exit 1
    fi
fi

echo -e "${BLUE}Configurações:${NC}"
echo "  Project ID: $PROJECT_ID"
echo "  Service Name: $SERVICE_NAME"
echo "  Region: $REGION"
echo "  Image: $IMAGE_NAME"
echo "  Memory: $MEMORY"
echo "  CPU: $CPU"
echo "  Min Instances: $MIN_INSTANCES"
echo "  Max Instances: $MAX_INSTANCES"
echo ""

read -p "Continuar com o deploy? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Deploy cancelado."
    exit 0
fi

# Configurar projeto
echo -e "${BLUE}1. Configurando projeto GCP...${NC}"
gcloud config set project $PROJECT_ID

# Habilitar APIs necessárias
echo -e "${BLUE}2. Habilitando APIs necessárias...${NC}"
gcloud services enable cloudbuild.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable containerregistry.googleapis.com

# Build da imagem usando Cloud Build
echo -e "${BLUE}3. Building imagem Docker com Cloud Build...${NC}"
gcloud builds submit --tag $IMAGE_NAME

# Deploy no Cloud Run
echo -e "${BLUE}4. Fazendo deploy no Cloud Run...${NC}"
gcloud run deploy $SERVICE_NAME \
  --image $IMAGE_NAME \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --memory $MEMORY \
  --cpu $CPU \
  --timeout ${TIMEOUT}s \
  --min-instances $MIN_INSTANCES \
  --max-instances $MAX_INSTANCES \
  --set-env-vars "ENVIRONMENT=production" \
  --set-env-vars "PORT=8080" \
  --set-env-vars "CAIXA_BASE_URL=https://servicebus2.caixa.gov.br/portaldeloterias/api" \
  --set-env-vars "CAIXA_TIMEOUT=15" \
  --set-env-vars "CAIXA_RATE_LIMIT=2" \
  --set-env-vars "CAIXA_CACHE_TTL=5" \
  --set-env-vars "CAIXA_MAX_RETRIES=3" \
  --set-env-vars "CAIXA_RETRY_BACKOFF=1000" \
  --set-env-vars "ASYNC_FORWARDING=false"

# Obter URL do serviço
echo ""
echo -e "${BLUE}5. Obtendo URL do serviço...${NC}"
SERVICE_URL=$(gcloud run services describe $SERVICE_NAME --platform managed --region $REGION --format 'value(status.url)')

echo ""
echo "========================================="
echo -e "${GREEN}✓ Deploy concluído com sucesso!${NC}"
echo "========================================="
echo ""
echo "URL do serviço: $SERVICE_URL"
echo ""
echo "Endpoints disponíveis:"
echo "  - Health: ${SERVICE_URL}/health"
echo "  - API: ${SERVICE_URL}/api/v1/lottery?gameType=megasena"
echo "  - Metrics: ${SERVICE_URL}/metrics"
echo ""
echo "Testar:"
echo "  curl ${SERVICE_URL}/health"
echo "  curl \"${SERVICE_URL}/api/v1/lottery?gameType=megasena\""
echo ""
echo "Ver logs:"
echo "  gcloud run logs read --service=$SERVICE_NAME --region=$REGION"
echo ""
echo "Gerenciar serviço:"
echo "  Console: https://console.cloud.google.com/run/detail/${REGION}/${SERVICE_NAME}"
echo ""
