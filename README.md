<p align="center">
    <img src="./.github/logo.png" width="200px">
</p>

<h1 align="center" style="font-weight: bold;">Desafio Técnico da V3</h1>

## ❤️ Bem vindos

Olá, tudo certo?

Seja bem vindo ao teste de seleção para novos colaboradores na V3!

Estamos honrados que você tenha chegado até aqui!

Prepare aquele ☕️, e venha conosco codar e se divertir!

## 📚 Desafios Disponíveis

Este repositório contém três desafios diferentes, cada um focado em uma área específica:

1. [Suporte Técnico](SUPPORT.md)
2. [Desafio Backend](CLOUD.md)
3. [Desafio Firmware](FIRMWARE.md)
   
### Endpoints

#### Endpoints para envio de dados

- `POST /telemetry/gyroscope` - Recebe dados do giroscópio
- `POST /telemetry/gps` - Recebe dados de GPS
- `POST /telemetry/photo` - Recebe fotos em formato base64

#### Endpoints para consulta de dados

- `GET /telemetry/gyroscope` - Retorna todos os dados do giroscópio
- `GET /telemetry/gps` - Retorna todos os dados de GPS
- `GET /telemetry/photo` - Retorna todos os dados de fotos

### Scripts

- `scripts/start.sh` - Inicia o container PostgreSQL e executa a aplicação
- `scripts/test_integration.sh` - Executa os testes de integração
- `scripts/build_and_run.sh` - Constrói e executa a aplicação em container Docker
- `test/run_tests.sh` - Executa todos os testes
- `scripts/start_postgres.sh` - Inicia o PostgreSQL localmente usando Docker
- `scripts/run_app.sh` - Compila e executa a aplicação localmente
- `scripts/run_local_tests.sh` - Executa os testes unitários e de integração localmente
- `scripts/docker_start.sh` - Inicia todos os serviços usando Docker Compose
- `scripts/test_docker.sh` - Inicia os serviços com Docker Compose e testa todos os endpoints

Documentação completa da API disponível em [docs/api.md](docs/api.md)

### Docker

A aplicação pode ser executada facilmente com Docker Compose, que configura tanto a API quanto o banco de dados PostgreSQL:

```bash
./scripts/docker_start.sh
```

Ou manualmente:

```bash
docker-compose up -d
```

A API estará disponível em http://localhost:8080

### CI/CD com GitHub Actions

O projeto está configurado com GitHub Actions para:

1. Executar testes unitários em cada push e pull request
2. Construir e publicar a imagem Docker no Docker Hub quando há push na branch main

Para configurar o CI/CD, adicione os seguintes secrets no repositório GitHub:

- `DOCKERHUB_USERNAME`: Seu nome de usuário do Docker Hub
- `DOCKERHUB_TOKEN`: Token de acesso do Docker Hub

## Poxa, outro teste?

Nós sabemos que os processos de seleção podem ser ingratos! Você investe um tempão e no final pode não ser aprovado!

Aqui, nós presamos pela **transparência**!

Este teste tem um **propósito** bastante simples:

> Nós queremos avaliar como você consegue transformar problemas em soluções através de código!

**🚨 IMPORTANTE!** Se você entende que já possui algum projeto pessoal, ou contribuição em um projeto _open-source_ que contemple conhecimentos equivalentes aos que existem neste desafio, então, basta submeter o repositório explicando essa correlação!

## 🚀 Bora nessa!

Este é um teste para analisarmos como você desempenha ao entender, traduzir, resolver e entregar um código que resolve um problema.

### Dicas

- Documente seu projeto;
- Faça perguntas sobre os pontos que não ficaram claros para você;
- Mostre a sua linha de raciocínio;
- Trabalhe bem o seu README.md;
  - Explique até onde implementou;
  - Como o projeto pode ser executado;
  - Como pode-se testar o projeto;

### Como você deverá desenvolver?

1. Faça um _fork_ deste projeto em seu GitHub pessoal;
2. Realize as implementações de acordo com cada um dos níveis;
3. Faça pequenos _commits_;
4. Depois de sentir que fez o seu máximo, faça um PR para o repositório original.

🚨 **IMPORTANTE!** Não significa que você precisa implementar **todos os níveis** para ser aprovado no processo! Faça até onde se sentir confortável.

## ⏰ Tempo para Entrega

Quanto antes você enviar, mais cuidado podemos ter na revisão do seu teste. Faça no seu tempo, mas mantenha a qualidade!

**Mas não desista! Envie até onde conseguir.**

Boa sorte! 🍀