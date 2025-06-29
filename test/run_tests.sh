#!/bin/bash

echo "Executando testes unitários..."
go test -v ./test/services/...

echo -e "\nExecutando testes de integração..."
go test -v ./test/repositories/...
