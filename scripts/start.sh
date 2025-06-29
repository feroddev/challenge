#!/bin/bash

echo "Iniciando ambiente de desenvolvimento..."
docker-compose up -d postgres

echo "Aguardando inicialização do PostgreSQL..."
sleep 5

echo "Executando a aplicação..."
go run cmd/app/main.go
