#!/bin/bash

# Script para testar a API localmente
# Author: Claude Code

# Cores
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

BASE_URL="http://localhost:8080"

echo "========================================="
echo "  Lottery Caixa Service - API Tests"
echo "========================================="
echo ""

# Verificar se o serviço está rodando
echo -e "${BLUE}1. Verificando se o serviço está rodando...${NC}"
if curl -s "${BASE_URL}/health" > /dev/null; then
    echo -e "${GREEN}✓ Serviço está rodando${NC}"
else
    echo -e "${RED}✗ Serviço não está rodando. Execute './start-local.sh' primeiro.${NC}"
    exit 1
fi
echo ""

# Health Check
echo -e "${BLUE}2. Health Check${NC}"
echo "GET ${BASE_URL}/health"
curl -s "${BASE_URL}/health" | jq '.' 2>/dev/null || curl -s "${BASE_URL}/health"
echo ""
echo ""

# Readiness Check
echo -e "${BLUE}3. Readiness Check${NC}"
echo "GET ${BASE_URL}/ready"
curl -s "${BASE_URL}/ready" | jq '.' 2>/dev/null || curl -s "${BASE_URL}/ready"
echo ""
echo ""

# Buscar resultado - Lotofácil
echo -e "${BLUE}4. Buscar Resultado - Lotofácil${NC}"
echo "GET ${BASE_URL}/?gameType=lotofacil"
curl -s "${BASE_URL}/?gameType=lotofacil" | jq '.' 2>/dev/null || curl -s "${BASE_URL}/?gameType=lotofacil"
echo ""
echo ""

# Buscar resultado - Mega Sena
echo -e "${BLUE}5. Buscar Resultado - Mega Sena${NC}"
echo "GET ${BASE_URL}/?gameType=megasena"
curl -s "${BASE_URL}/?gameType=megasena" | jq '.' 2>/dev/null || curl -s "${BASE_URL}/?gameType=megasena"
echo ""
echo ""

# Webhook - Enviar dados
echo -e "${BLUE}6. Webhook - Enviar Dados de Sorteio${NC}"
echo "POST ${BASE_URL}/api/v1/lottery/webhook"
curl -s -X POST "${BASE_URL}/api/v1/lottery/webhook" \
  -H "Content-Type: application/json" \
  -d '{
    "gameType": "lotofacil",
    "drawNumber": 2024001,
    "drawDate": "2024-01-15",
    "numbers": [1, 5, 8, 12, 15, 18, 22, 25],
    "timestamp": "2024-01-15T10:30:00Z"
  }' | jq '.' 2>/dev/null || curl -s -X POST "${BASE_URL}/api/v1/lottery/webhook" \
  -H "Content-Type: application/json" \
  -d '{
    "gameType": "lotofacil",
    "drawNumber": 2024001,
    "drawDate": "2024-01-15",
    "numbers": [1, 5, 8, 12, 15, 18, 22, 25],
    "timestamp": "2024-01-15T10:30:00Z"
  }'
echo ""
echo ""

# Métricas (primeiras 20 linhas)
echo -e "${BLUE}7. Métricas Prometheus (primeiras 20 linhas)${NC}"
echo "GET ${BASE_URL}/metrics"
curl -s "${BASE_URL}/metrics" | head -20
echo "..."
echo ""

# Verificar downstream (se estiver rodando)
echo -e "${BLUE}8. Verificar Downstream Processor${NC}"
echo "GET http://localhost:8081/health"
if curl -s "http://localhost:8081/health" > /dev/null 2>&1; then
    curl -s "http://localhost:8081/health" | jq '.' 2>/dev/null || curl -s "http://localhost:8081/health"
    echo -e "${GREEN}✓ Downstream está rodando${NC}"
else
    echo -e "${YELLOW}⚠️  Downstream não está rodando (opcional)${NC}"
fi
echo ""

echo "========================================="
echo -e "${GREEN}✓ Testes concluídos!${NC}"
echo "========================================="
echo ""
echo "Dicas:"
echo "  - Use 'jq' para formatar JSON: apt-get install jq"
echo "  - Verifique logs do serviço no terminal onde está rodando"
echo "  - Métricas completas: curl http://localhost:8080/metrics"
echo "  - Prometheus: http://localhost:9090 (se usando Docker Compose)"
echo "  - Grafana: http://localhost:3000 (se usando Docker Compose)"
echo ""
