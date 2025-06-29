# Sistema de Mensagens NATS

## Visão Geral

O sistema de mensagens NATS foi implementado para permitir o processamento assíncrono de dados de telemetria. Isso proporciona maior escalabilidade, resiliência e desacoplamento entre a API que recebe os dados e os serviços que os processam.

## Arquitetura

### Produtor

O produtor é responsável por publicar mensagens nos tópicos NATS quando a API recebe dados de telemetria. Ele está integrado ao serviço de telemetria e publica mensagens nos seguintes tópicos:

- `gyroscope.telemetry`: Dados do giroscópio
- `gps.telemetry`: Dados de GPS
- `photo.telemetry`: Dados de fotos

O produtor implementa um mecanismo de retry para garantir que as mensagens sejam publicadas mesmo em caso de falhas temporárias na conexão com o servidor NATS. Após esgotar as tentativas de retry, as mensagens são enviadas para um tópico de dead-letter para análise posterior.

### Consumidores

Os consumidores são serviços independentes que assinam os tópicos NATS e processam as mensagens recebidas:

1. **Consumidor de Giroscópio**:
   - Assina o tópico `gyroscope.telemetry`
   - Desserializa os dados para o modelo `Gyroscope`
   - Persiste os dados no banco de dados PostgreSQL
   - Registra logs estruturados do processamento

2. **Consumidor de GPS**:
   - Assina o tópico `gps.telemetry`
   - Desserializa os dados para o modelo `GPS`
   - Persiste os dados no banco de dados PostgreSQL
   - Registra logs estruturados do processamento

3. **Consumidor de Fotos**:
   - Assina o tópico `photo.telemetry`
   - Desserializa os dados para o modelo `Photo`
   - Persiste os dados no banco de dados PostgreSQL
   - Inicia o processamento de reconhecimento facial usando AWS Rekognition
   - Utiliza cache Redis para otimizar consultas e evitar chamadas repetidas
   - Registra logs estruturados do processamento

### Mecanismos de Resiliência

O sistema implementa vários mecanismos para garantir a resiliência no processamento de mensagens:

1. **Retry Automático**:
   - Tanto o produtor quanto os consumidores implementam mecanismos de retry
   - O número de tentativas e o intervalo entre elas são configuráveis via variáveis de ambiente

2. **Dead-Letter Queue**:
   - Mensagens que não puderam ser processadas após todas as tentativas de retry são enviadas para um tópico de dead-letter
   - Isso permite análise posterior e eventual reprocessamento manual

3. **Logs Estruturados**:
   - Todo o ciclo de vida das mensagens é registrado com logs estruturados usando Zap
   - Isso facilita o monitoramento e diagnóstico de problemas

## Configuração

A configuração do sistema de mensagens NATS é feita através de variáveis de ambiente:

```
NATS_URL=nats://localhost:4222
NATS_GYROSCOPE_TOPIC=gyroscope.telemetry
NATS_GPS_TOPIC=gps.telemetry
NATS_PHOTO_TOPIC=photo.telemetry
NATS_RETRY_ATTEMPTS=3
NATS_RETRY_DELAY=5
NATS_DEAD_LETTER_TOPIC=telemetry.deadletter
```

## Execução

Os consumidores podem ser executados de várias formas:

1. **Docker Compose**:
   ```bash
   docker-compose up -d
   ```

2. **Script de Execução**:
   ```bash
   ./scripts/run_services.sh
   ```

3. **Execução Manual**:
   ```bash
   # Compilar os binários
   go build -o bin/gyroscope-consumer cmd/consumers/gyroscope/main.go
   go build -o bin/gps-consumer cmd/consumers/gps/main.go
   go build -o bin/photo-consumer cmd/consumers/photo/main.go

   # Executar os consumidores
   ./bin/gyroscope-consumer &
   ./bin/gps-consumer &
   ./bin/photo-consumer &
   ```

## Monitoramento

O servidor NATS expõe uma interface de monitoramento na porta 8222. Você pode acessá-la em http://localhost:8222 quando o servidor estiver em execução.

## Considerações de Escalabilidade

O sistema foi projetado para ser escalável:

1. Cada consumidor pode ser escalado horizontalmente para processar mais mensagens
2. Os consumidores são independentes, permitindo escalabilidade seletiva
3. O NATS suporta balanceamento de carga entre múltiplas instâncias de consumidores
