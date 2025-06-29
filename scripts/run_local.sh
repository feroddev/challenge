#!/bin/bash

# Cores para saída
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Iniciando serviços da API de Telemetria (Modo Local)${NC}"

# Configurar variáveis de ambiente para serviços Docker
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=telemetry
export DB_SSLMODE=disable

export REDIS_ADDR=localhost:6379
export REDIS_PASSWORD=""

export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=sua_access_key
export AWS_SECRET_ACCESS_KEY=sua_secret_key

export NATS_URL=nats://localhost:4222
export NATS_GYROSCOPE_TOPIC=gyroscope.telemetry
export NATS_GPS_TOPIC=gps.telemetry
export NATS_PHOTO_TOPIC=photo.telemetry
export NATS_RETRY_ATTEMPTS=3
export NATS_RETRY_DELAY=5
export NATS_DEAD_LETTER_TOPIC=telemetry.deadletter

export PORT=8081
export ENV=development

# Compila os binários
echo -e "${YELLOW}Compilando binários...${NC}"
mkdir -p bin

echo -e "${YELLOW}Compilando API...${NC}"
go build -o bin/api cmd/app/main.go

echo -e "${YELLOW}Compilando consumidor de giroscópio...${NC}"
go build -o bin/gyroscope-consumer cmd/consumers/gyroscope/main.go

echo -e "${YELLOW}Compilando consumidor de GPS...${NC}"
go build -o bin/gps-consumer cmd/consumers/gps/main.go

echo -e "${YELLOW}Compilando consumidor de fotos...${NC}"
go build -o bin/photo-consumer cmd/consumers/photo/main.go

# Inicia os serviços em background
echo -e "${YELLOW}Iniciando serviços...${NC}"

echo -e "${GREEN}Iniciando API...${NC}"
./bin/api &
API_PID=$!

echo -e "${GREEN}Iniciando consumidor de giroscópio...${NC}"
./bin/gyroscope-consumer &
GYRO_PID=$!

echo -e "${GREEN}Iniciando consumidor de GPS...${NC}"
./bin/gps-consumer &
GPS_PID=$!

echo -e "${GREEN}Iniciando consumidor de fotos...${NC}"
./bin/photo-consumer &
PHOTO_PID=$!

# Função para encerrar todos os processos
function cleanup {
  echo -e "${YELLOW}Encerrando serviços...${NC}"
  kill $API_PID $GYRO_PID $GPS_PID $PHOTO_PID
  wait
  echo -e "${GREEN}Serviços encerrados.${NC}"
}

# Registra trap para SIGINT e SIGTERM
trap cleanup SIGINT SIGTERM

echo -e "${GREEN}Todos os serviços estão rodando!${NC}"
echo -e "${GREEN}API: http://localhost:${PORT:-8080}${NC}"
echo -e "${YELLOW}Pressione Ctrl+C para encerrar todos os serviços${NC}"

# Aguarda até que o usuário pressione Ctrl+C
wait
