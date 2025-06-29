#!/bin/bash

echo "Iniciando banco de dados de teste..."
docker-compose up -d postgres_test

echo "Aguardando inicialização do PostgreSQL de teste..."
sleep 5

echo "Executando testes de integração..."
go test -v ./test/repositories/...
