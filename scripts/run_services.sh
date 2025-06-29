#!/bin/bash

# Cores para saída
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Iniciando serviços da API de Telemetria${NC}"

# Carrega variáveis de ambiente
if [ -f .env ]; then
  echo -e "${GREEN}Carregando variáveis de ambiente do arquivo .env${NC}"
  export $(cat .env | grep -v '^#' | xargs)
else
  echo -e "${YELLOW}Arquivo .env não encontrado, usando variáveis de ambiente do sistema${NC}"
fi

# Verifica se o Docker está rodando
echo -e "${YELLOW}Verificando serviços Docker...${NC}"
if ! docker info > /dev/null 2>&1; then
  echo -e "${RED}Docker não está rodando. Por favor, inicie o Docker e tente novamente.${NC}"
  exit 1
fi

# Inicia os serviços de dependência (PostgreSQL, Redis e NATS)
echo -e "${YELLOW}Iniciando serviços de dependência (PostgreSQL, Redis e NATS)...${NC}"
docker-compose up -d postgres redis nats

# Aguarda os serviços estarem prontos
echo -e "${YELLOW}Aguardando serviços estarem prontos...${NC}"
sleep 5

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
