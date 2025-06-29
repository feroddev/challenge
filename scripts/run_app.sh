#!/bin/bash

echo "Compilando a aplicação..."
go build -o ./bin/app ./cmd/app

echo "Executando a aplicação..."
./bin/app
