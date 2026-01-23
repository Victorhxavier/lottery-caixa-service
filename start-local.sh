#!/bin/bash

# Script para inicializar o projeto localmente
# Author: Claude Code
# Description: Inicia o serviço de loteria e o processador downstream

set -e

echo "========================================="
echo "  Lottery Caixa Service - Local Setup"
echo "========================================="
echo ""

# Cores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Verificar se Go está instalado
if ! command -v go &> /dev/null; then
    echo -e "${YELLOW}⚠️  Go não está instalado. Por favor, instale Go 1.21 ou superior.${NC}"
    exit 1
fi

echo -e "${BLUE}✓ Go version:${NC} $(go version)"
echo ""

# Verificar se .env existe
if [ ! -f ".env" ]; then
    echo -e "${YELLOW}⚠️  Arquivo .env não encontrado. Criando a partir de .env.example...${NC}"
    cp .env.example .env
    echo -e "${GREEN}✓ Arquivo .env criado${NC}"
fi

# Instalar dependências
echo -e "${BLUE}📦 Baixando dependências...${NC}"
go mod download
go mod tidy
echo -e "${GREEN}✓ Dependências instaladas${NC}"
echo ""

# Perguntar qual modo de execução
echo "Selecione o modo de execução:"
echo "  1) Apenas o serviço principal (porta 8080)"
echo "  2) Serviço principal + Processador downstream (portas 8080 e 8081)"
echo "  3) Ambiente completo com Docker Compose (recomendado)"
echo ""
read -p "Escolha uma opção [1-3]: " choice

case $choice in
    1)
        echo ""
        echo -e "${BLUE}🚀 Iniciando serviço principal na porta 8080...${NC}"
        echo -e "${YELLOW}⚠️  Nota: O downstream não estará disponível. Algumas funcionalidades podem não funcionar.${NC}"
        echo ""
        go run ./cmd/main.go
        ;;
    2)
        echo ""
        echo -e "${BLUE}🚀 Iniciando serviços em terminais separados...${NC}"
        echo ""

        # Iniciar processador downstream em background
        echo -e "${BLUE}📡 Iniciando processador downstream (porta 8081)...${NC}"
        PORT=8081 go run ./examples/processor.go &
        PROCESSOR_PID=$!

        sleep 2

        # Iniciar serviço principal
        echo -e "${BLUE}🎰 Iniciando serviço principal (porta 8080)...${NC}"
        echo ""
        echo -e "${GREEN}✓ Serviços iniciados!${NC}"
        echo ""
        echo "Endpoints disponíveis:"
        echo "  - API Principal: http://localhost:8080"
        echo "  - Health Check:  http://localhost:8080/health"
        echo "  - Processador:   http://localhost:8081/health"
        echo ""
        echo -e "${YELLOW}Pressione Ctrl+C para parar os serviços${NC}"
        echo ""

        # Trap para limpar processos ao sair
        trap "echo ''; echo 'Parando serviços...'; kill $PROCESSOR_PID 2>/dev/null; exit" INT TERM

        go run ./cmd/main.go
        ;;
    3)
        echo ""
        echo -e "${BLUE}🐳 Iniciando ambiente com Docker Compose...${NC}"
        echo ""

        # Verificar se Docker está instalado
        if ! command -v docker &> /dev/null; then
            echo -e "${YELLOW}⚠️  Docker não está instalado. Por favor, instale Docker primeiro.${NC}"
            exit 1
        fi

        # Verificar se docker-compose está instalado
        if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
            echo -e "${YELLOW}⚠️  Docker Compose não está instalado.${NC}"
            exit 1
        fi

        echo "Iniciando todos os serviços (API, Processor, Prometheus, Grafana)..."
        docker-compose up --build
        ;;
    *)
        echo -e "${YELLOW}⚠️  Opção inválida${NC}"
        exit 1
        ;;
esac
