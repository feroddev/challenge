# Reconhecimento Facial com AWS Rekognition

## Visão Geral

A API de telemetria agora possui integração com AWS Rekognition para comparar fotos enviadas com fotos anteriores do mesmo dispositivo. O sistema utiliza Redis para cache e registra logs estruturados com Zap.

## Funcionalidades

- Comparação facial entre fotos do mesmo dispositivo
- Cache de resultados com Redis
- Logs estruturados com Zap
- Processamento assíncrono para não bloquear a API

## Configuração

### Variáveis de Ambiente

```
# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=

# AWS
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=sua_access_key
AWS_SECRET_ACCESS_KEY=sua_secret_key
```

### Docker Compose

O arquivo `docker-compose.yml` foi atualizado para incluir o serviço Redis:

```yaml
redis:
  image: redis:7
  container_name: telemetry_redis
  ports:
    - "6379:6379"
  networks:
    - telemetry_network
  volumes:
    - redis_data:/data
```

## Como Executar

### Usando Docker

```bash
# Configure suas credenciais AWS
export AWS_ACCESS_KEY_ID=sua_access_key
export AWS_SECRET_ACCESS_KEY=sua_secret_key

# Inicie os serviços
docker-compose up -d
```

### Usando Script

```bash
# Configure suas credenciais AWS
export AWS_ACCESS_KEY_ID=sua_access_key
export AWS_SECRET_ACCESS_KEY=sua_secret_key

# Execute o script
./scripts/run_with_aws.sh
```

## Testando o Reconhecimento Facial

Para testar a funcionalidade de reconhecimento facial, utilize o script `test_facial_recognition.sh`:

```bash
# Certifique-se que a API está em execução
./scripts/test_facial_recognition.sh foto1.jpg foto2.jpg foto3.jpg
```

Este script enviará as fotos para a API e exibirá os resultados do reconhecimento facial.

## Fluxo de Processamento

1. A foto é recebida pela API e salva no banco de dados
2. O sistema busca fotos anteriores do mesmo dispositivo (usando cache se disponível)
3. A nova foto é comparada com as fotos anteriores usando AWS Rekognition
4. O resultado da comparação é salvo no banco de dados e no cache
5. O endpoint GET retorna os resultados com informações de reconhecimento

## Estrutura de Logs

Os logs são estruturados usando Zap e incluem informações como:

- Método HTTP e path da requisição
- Status da resposta
- Latência da requisição
- IP do cliente
- Detalhes de operações de cache
- Resultados de comparação facial
- Erros e informações de depuração
