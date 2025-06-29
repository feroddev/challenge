#!/bin/bash

# Script para testar o reconhecimento facial da API

# Verifica se curl está instalado
if ! command -v curl &> /dev/null; then
    echo "curl não está instalado. Por favor instale-o primeiro."
    exit 1
fi

# Verifica se base64 está instalado
if ! command -v base64 &> /dev/null; then
    echo "base64 não está instalado. Por favor instale-o primeiro."
    exit 1
fi

# Configurações
API_URL="http://localhost:8080"
DEVICE_ID="teste123"
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Função para enviar uma foto para a API
enviar_foto() {
    local arquivo=$1
    local device_id=$2
    
    # Converte a imagem para base64
    BASE64_IMG=$(base64 -w 0 $arquivo)
    
    # Cria o payload JSON
    JSON_PAYLOAD=$(cat <<EOF
{
  "photo": "$BASE64_IMG",
  "timestamp": "$TIMESTAMP",
  "device_id": "$device_id"
}
EOF
)
    
    # Envia a requisição POST
    echo "Enviando foto $arquivo para o device_id $device_id..."
    curl -s -X POST "$API_URL/telemetry/photo" \
        -H "Content-Type: application/json" \
        -d "$JSON_PAYLOAD"
    
    echo -e "\n"
}

# Função para buscar dados de fotos
buscar_fotos() {
    echo "Buscando dados de fotos..."
    curl -s "$API_URL/telemetry/photo" | jq '.'
    echo -e "\n"
}

# Testa o fluxo completo
echo "Iniciando teste de reconhecimento facial"
echo "========================================"

# Verifica se foram fornecidas as imagens para teste
if [ $# -lt 2 ]; then
    echo "Uso: $0 <imagem1> <imagem2> [imagem3...]"
    echo "Exemplo: $0 foto1.jpg foto2.jpg foto3.jpg"
    exit 1
fi

# Envia cada imagem fornecida
for img in "$@"; do
    if [ ! -f "$img" ]; then
        echo "Arquivo $img não encontrado!"
        continue
    fi
    
    enviar_foto "$img" "$DEVICE_ID"
    
    # Aguarda um pouco para o processamento assíncrono
    echo "Aguardando processamento..."
    sleep 2
done

# Busca os resultados
buscar_fotos

echo "Teste concluído!"
