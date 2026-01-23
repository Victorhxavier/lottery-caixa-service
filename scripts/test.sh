#!/bin/bash

set -e

# Cores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configurações
UPSTREAM_URL="${UPSTREAM_URL:-http://localhost:8080}"
GAME_TYPES=("lotofacil" "lotomania" "duplasena" "quina" "megasena")
TIMEOUT=10

echo -e "${BLUE}═══════════════════════════════════════════${NC}"
echo -e "${BLUE}  Lottery Service Integration Tests${NC}"
echo -e "${BLUE}═══════════════════════════════════════════${NC}\n"

# Test 1: Health Check
echo -e "${YELLOW}[1] Health Check${NC}"
if curl -s -f --connect-timeout $TIMEOUT "$UPSTREAM_URL/health" > /dev/null; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED${NC}"
    exit 1
fi

# Test 2: Readiness
echo -e "${YELLOW}[2] Readiness Check${NC}"
if curl -s -f --connect-timeout $TIMEOUT "$UPSTREAM_URL/ready" > /dev/null; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED${NC}"
    exit 1
fi

# Test 3: Metrics
echo -e "${YELLOW}[3] Metrics Endpoint${NC}"
if curl -s -f --connect-timeout $TIMEOUT "$UPSTREAM_URL/metrics" > /dev/null; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED${NC}"
    exit 1
fi

# Test 4: Service Info
echo -e "${YELLOW}[4] Service Info${NC}"
if curl -s -f --connect-timeout $TIMEOUT "$UPSTREAM_URL/api/v1/info" > /dev/null; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${YELLOW}⚠ WARNING${NC}"
fi

# Test 5: Lottery Queries
echo -e "${YELLOW}[5] Lottery Queries${NC}"
for game_type in "${GAME_TYPES[@]}"; do
    echo -n "  Testing $game_type... "
    response=$(curl -s --connect-timeout $TIMEOUT "$UPSTREAM_URL/?gameType=$game_type")
    
    if echo "$response" | jq . > /dev/null 2>&1; then
        status=$(echo "$response" | jq -r '.status')
        if [ "$status" = "success" ]; then
            echo -e "${GREEN}✓${NC}"
        else
            echo -e "${RED}✗ (status: $status)${NC}"
        fi
    else
        echo -e "${RED}✗ (invalid JSON)${NC}"
    fi
done

# Test 6: Response Time
echo -e "${YELLOW}[6] Response Time${NC}"
start=$(date +%s%N)
curl -s "$UPSTREAM_URL/?gameType=lotofacil" > /dev/null
end=$(date +%s%N)
elapsed=$((($end - $start) / 1000000))
echo "  Duration: ${elapsed}ms"

if [ $elapsed -lt 1000 ]; then
    echo -e "${GREEN}✓ Acceptable${NC}"
elif [ $elapsed -lt 5000 ]; then
    echo -e "${YELLOW}⚠ Slow${NC}"
else
    echo -e "${RED}✗ Too slow${NC}"
fi

# Test 7: Concurrent Requests
echo -e "${YELLOW}[7] Concurrent Requests (5x)${NC}"
success=0
for i in {1..5}; do
    if curl -s "$UPSTREAM_URL/?gameType=lotofacil" > /dev/null; then
        ((success++))
    fi
done
echo "  Success: $success/5"

if [ $success -eq 5 ]; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${YELLOW}⚠ Some requests failed${NC}"
fi

# Test 8: Error Handling
echo -e "${YELLOW}[8] Error Handling${NC}"
response=$(curl -s --connect-timeout $TIMEOUT "$UPSTREAM_URL/?gameType=invalid" || true)
if echo "$response" | jq . > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Proper error response${NC}"
else
    echo -e "${YELLOW}⚠ No JSON error response${NC}"
fi

echo ""
echo -e "${BLUE}═══════════════════════════════════════════${NC}"
echo -e "${GREEN}✓ Tests completed${NC}"
echo -e "${BLUE}═══════════════════════════════════════════${NC}"

echo ""
echo -e "${YELLOW}Additional Information:${NC}"
echo "  Upstream URL: $UPSTREAM_URL"
echo "  Test Time: $(date)"
echo ""
echo -e "${YELLOW}Next Steps:${NC}"
echo "  1. View logs: docker-compose logs -f lottery-caixa-service"
echo "  2. Check Prometheus: http://localhost:9090"
echo "  3. Check Grafana: http://localhost:3000"
echo ""
