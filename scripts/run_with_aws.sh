#!/bin/bash

# Verifica se as variáveis de ambiente AWS estão definidas
if [ -z "$AWS_ACCESS_KEY_ID" ] || [ -z "$AWS_SECRET_ACCESS_KEY" ]; then
  echo "Erro: Variáveis de ambiente AWS_ACCESS_KEY_ID e AWS_SECRET_ACCESS_KEY devem estar definidas."
  echo "Execute: export AWS_ACCESS_KEY_ID=sua_chave_acesso AWS_SECRET_ACCESS_KEY=sua_chave_secreta"
  exit 1
fi

# Inicia os serviços com Docker Compose
docker-compose up -d postgres redis

# Aguarda os serviços estarem prontos
echo "Aguardando serviços iniciarem..."
sleep 5

# Executa a aplicação com as variáveis de ambiente necessárias
export REDIS_ADDR=localhost:6379
export REDIS_PASSWORD=""
export AWS_REGION=us-east-1

# Executa a aplicação
cd ..
go run cmd/app/main.go
