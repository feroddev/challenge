#!/bin/bash

echo "Iniciando PostgreSQL para desenvolvimento..."
docker run --name telemetry_postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=telemetry -p 5432:5432 -d postgres:14

echo "Aguardando inicialização do PostgreSQL..."
sleep 5

echo "Criando banco de dados de teste..."
docker exec telemetry_postgres psql -U postgres -c "CREATE DATABASE telemetry_test;"

echo "PostgreSQL iniciado com sucesso!"
echo "- Banco principal: telemetry"
echo "- Banco de teste: telemetry_test"
echo "- Usuário: postgres"
echo "- Senha: postgres"
echo "- Porta: 5432"
