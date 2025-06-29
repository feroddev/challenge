# Arquitetura do Sistema de Telemetria

## Diagrama de Arquitetura

```
+------------------+     +------------------+     +------------------+
|                  |     |                  |     |                  |
|  Dispositivos    |---->|  API REST        |---->|  PostgreSQL      |
|  Android         |     |  (Go/Gin)        |     |  (Persistência)  |
|                  |     |                  |     |                  |
+------------------+     +--------+---------+     +------------------+
                                  |
                                  |
                                  v
+------------------+     +------------------+     +------------------+
|                  |     |                  |     |                  |
|  Consumidores    |<----|  NATS           |<----|  Métricas        |
|  (Processamento) |     |  (Mensageria)   |     |  (Prometheus)    |
|                  |     |                  |     |                  |
+------------------+     +--------+---------+     +------------------+
                                  |
                                  |
                                  v
+------------------+     +------------------+     +------------------+
|                  |     |                  |     |                  |
|  AWS Rekognition |<--->|  Redis Cache     |<--->|  OpenTelemetry  |
|  (Reconhecimento)|     |  (Performance)   |     |  (Tracing)       |
|                  |     |                  |     |                  |
+------------------+     +------------------+     +------------------+
```

## Componentes Principais

### API REST (Go/Gin)
- Recebe dados de telemetria dos dispositivos Android
- Implementa endpoints para giroscópio, GPS e fotos
- Documentação Swagger integrada
- Métricas Prometheus para monitoramento
- Tracing com OpenTelemetry

### Persistência (PostgreSQL)
- Armazena todos os dados de telemetria
- Tabelas para giroscópio, GPS e fotos
- Implementado com GORM

### Mensageria (NATS)
- Sistema de mensagens pub/sub
- Garante processamento assíncrono
- Implementa mecanismo de retry e dead-letter

### Cache (Redis)
- Armazena resultados de reconhecimento facial
- Melhora performance reduzindo chamadas ao banco e AWS

### Reconhecimento Facial (AWS Rekognition)
- Compara fotos do mesmo dispositivo
- Detecta similaridade entre imagens

### Monitoramento e Observabilidade
- Métricas com Prometheus
- Tracing com OpenTelemetry
- Logs estruturados com Zap

### Consumidores
- Processam mensagens de telemetria assincronamente
- Implementam lógica específica para cada tipo de dado
- Executados em containers separados

## Fluxo de Dados

1. Dispositivo Android envia dados para a API REST
2. API valida e persiste os dados no PostgreSQL
3. API publica mensagem no NATS
4. Consumidores processam as mensagens
5. Para fotos, o reconhecimento facial é executado
6. Resultados são armazenados no Redis e PostgreSQL
7. Métricas e traces são coletados em todo o processo
