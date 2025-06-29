#!/bin/bash

echo "Construindo imagem Docker..."
docker build -t telemetry-api .

echo "Executando container..."
docker run -p 8080:8080 --network telemetry_network -e DB_HOST=postgres -e DB_PORT=5432 -e DB_USER=postgres -e DB_PASSWORD=postgres -e DB_NAME=telemetry -e DB_SSLMODE=disable telemetry-api
