## Rate Limiter em Go com Redis

Este projeto implementa um rate limiter em Go (Golang) que utiliza o Redis como armazenamento para controlar o número de requisições HTTP permitidas por segundo, com base no endereço IP ou em um token de acesso.

### Funcionalidades

- **Limitação por IP:** Limita o número de requisições por segundo de um determinado endereço IP.
- **Limitação por Token:** Limita o número de requisições por segundo para um token de acesso específico (enviado no cabeçalho `API_KEY`).
- **Bloqueio Temporário:** Bloqueia o IP ou token por um período de tempo configurável (5 minutos por padrão) se o limite for excedido.
- **Cabeçalho Retry-After:** Inclui o cabeçalho `Retry-After` nas respostas bloqueadas, informando ao cliente quanto tempo ele deve esperar antes de fazer outra requisição.
- **Armazenamento em Redis:** Utiliza o Redis como banco de dados para armazenar os contadores de requisições e os tempos de bloqueio.
- **Estratégia Flexível:** Permite a implementação de diferentes estratégias de armazenamento (além do Redis) através de uma interface.
- **Middleware:** Fornece um middleware para integrar o rate limiter ao seu servidor web.

### Configuração

As configurações do rate limiter são definidas no arquivo `app.env`:

- `LIMIT_BY`: Define o tipo de limitação ("ip" ou "token").
- `MAX_REQUESTS_PER_SECOND_IP`: Define o número máximo de requisições por IP permitidas por segundo.
- `MAX_REQUESTS_PER_SECOND_TOKEN`: Define o número máximo de requisições por token permitidas por segundo.
- `BLOCK_DURATION_SECONDS`: Define a duração do bloqueio em segundos (5 minutos por padrão).
- `REDIS_ADDR`: Define o endereço do servidor Redis.

### Estrutura do Projeto

```
rate-limiter/
├── cmd/
│   └── main.go (Servidor web)
├── internal/
│   ├── config/
│   │   └── config.go (Carregamento de configurações)
│   ├── limiter/
│   │   ├── limiter.go (Interface)
│   │   ├── iplimiter.go (Implementação por IP)
│   │   ├── tokenlimiter.go (Implementação por token)
│   │   └── strategy/
│   │       └── redis.go (Estratégia Redis)
│   └── middleware/
│       └── middleware.go
├── Dockerfile
├── docker-compose.yml
└── app.env
```

### Como Usar

1. **Clone o Repositório:**

   ```bash
   git clone https://github.com/Sanpeta/rate-limiter-pos-go-expert.git
   ```

2. **Renomear o arquivo `example.env` para `app.env` e configurar as variáveis:**

   ```bash
   LIMIT_BY=ip OU token
   MAX_REQUESTS_PER_SECOND_IP=5
   MAX_REQUESTS_PER_SECOND_TOKEN=10
   BLOCK_DURATION_SECONDS=300s
   REDIS_ADDR=redis:6379
   ```

3. **Inicie e Execute:**

   ```bash
   docker-compose up
   ```

4. **Teste:**

   - Envie requisições HTTP `GET` para `http://localhost:8080/` usando `curl`, `Postman` ou outra ferramenta.
   - Você pode utilizar o arquivo `test_rate_limiter.sh` e modificar o número de interações que deseja realizar.


**Observações:**

- Certifique-se de que o Docker e o Docker Compose estejam instalados em sua máquina.
- O servidor web responderá na porta 8080.
- O Redis será iniciado automaticamente pelo Docker Compose.
- As configurações podem ser ajustadas no arquivo `app.env`.
- O script `test_rate_limiter.sh` pode ser utilizado para testar o funcionamento do rate limiter.
- Ajuste o número de requisições no script para simular diferentes cenários de uso.


**Contribuição:**

Sinta-se à vontade para contribuir com melhorias, correções de bugs ou novas funcionalidades!
