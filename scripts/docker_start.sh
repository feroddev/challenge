#!/bin/bash

echo "Iniciando serviços com Docker Compose..."
docker-compose up -d

echo "Serviços iniciados:"
echo "- API: http://localhost:8080"
echo "- PostgreSQL: localhost:5432"

echo "Para verificar os logs da API:"
echo "docker-compose logs -f app"
