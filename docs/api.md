# Documentação da API de Telemetria

## Visão Geral

Esta API foi desenvolvida para receber e armazenar dados de telemetria de dispositivos móveis, incluindo dados de giroscópio, GPS e fotos. A API utiliza PostgreSQL para persistência de dados e segue uma arquitetura em camadas com padrão repositório.

## Tecnologias Utilizadas

- Golang com Gin Framework
- PostgreSQL com GORM
- Docker para ambiente de desenvolvimento
- Testes unitários e de integração

## Endpoints

### Endpoints para Envio de Dados

#### POST /telemetry/gyroscope

Recebe dados do giroscópio de um dispositivo.

**Corpo da Requisição:**
```json
{
  "x": 10.5,
  "y": -5.2,
  "z": 3.7,
  "timestamp": "2025-06-28T20:30:00Z",
  "device_id": "abc123"
}
```

**Resposta de Sucesso:**
```json
{
  "message": "Dados do giroscópio recebidos com sucesso"
}
```

#### POST /telemetry/gps

Recebe dados de GPS de um dispositivo.

**Corpo da Requisição:**
```json
{
  "latitude": -23.5505,
  "longitude": -46.6333,
  "timestamp": "2025-06-28T20:30:00Z",
  "device_id": "abc123"
}
```

**Resposta de Sucesso:**
```json
{
  "message": "Dados do GPS recebidos com sucesso"
}
```

#### POST /telemetry/photo

Recebe fotos em formato base64 de um dispositivo.

**Corpo da Requisição:**
```json
{
  "photo": "base64_encoded_string",
  "timestamp": "2025-06-28T20:30:00Z",
  "device_id": "abc123"
}
```

**Resposta de Sucesso:**
```json
{
  "message": "Dados da foto recebidos com sucesso"
}
```

### Endpoints para Consulta de Dados

#### GET /telemetry/gyroscope

Retorna todos os dados de giroscópio armazenados.

**Resposta de Sucesso:**
```json
[
  {
    "id": 1,
    "x": 10.5,
    "y": -5.2,
    "z": 3.7,
    "timestamp": "2025-06-28T20:30:00Z",
    "device_id": "abc123",
    "created_at": "2025-06-28T20:30:05Z",
    "updated_at": "2025-06-28T20:30:05Z"
  },
  {
    "id": 2,
    "x": 11.2,
    "y": -4.8,
    "z": 3.9,
    "timestamp": "2025-06-28T20:31:00Z",
    "device_id": "abc123",
    "created_at": "2025-06-28T20:31:05Z",
    "updated_at": "2025-06-28T20:31:05Z"
  }
]
```

#### GET /telemetry/gps

Retorna todos os dados de GPS armazenados.

**Resposta de Sucesso:**
```json
[
  {
    "id": 1,
    "latitude": -23.5505,
    "longitude": -46.6333,
    "timestamp": "2025-06-28T20:30:00Z",
    "device_id": "abc123",
    "created_at": "2025-06-28T20:30:05Z",
    "updated_at": "2025-06-28T20:30:05Z"
  },
  {
    "id": 2,
    "latitude": -23.5506,
    "longitude": -46.6334,
    "timestamp": "2025-06-28T20:31:00Z",
    "device_id": "abc123",
    "created_at": "2025-06-28T20:31:05Z",
    "updated_at": "2025-06-28T20:31:05Z"
  }
]
```

#### GET /telemetry/photo

Retorna todos os dados de fotos armazenados.

**Resposta de Sucesso:**
```json
[
  {
    "id": 1,
    "photo": "base64_encoded_string",
    "timestamp": "2025-06-28T20:30:00Z",
    "device_id": "abc123",
    "created_at": "2025-06-28T20:30:05Z",
    "updated_at": "2025-06-28T20:30:05Z"
  },
  {
    "id": 2,
    "photo": "base64_encoded_string",
    "timestamp": "2025-06-28T20:31:00Z",
    "device_id": "abc123",
    "created_at": "2025-06-28T20:31:05Z",
    "updated_at": "2025-06-28T20:31:05Z"
  }
]
```

## Códigos de Status

- `201 Created`: Dados recebidos com sucesso
- `200 OK`: Consulta realizada com sucesso
- `400 Bad Request`: Dados inválidos
- `500 Internal Server Error`: Erro no servidor
