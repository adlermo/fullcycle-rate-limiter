# Rate Limiter - Go Expert Challenge

Implementação de um Rate Limiter em Go utilizando Redis para persistência e Docker para execução.

## Objetivo

Controlar o número de requisições por:

* IP
* Token (`API_KEY`)

Respeitando a regra:

> Configurações de Token possuem prioridade sobre configurações de IP.

---

## Funcionalidades

* Limitação por IP
* Limitação por Token
* Precedência Token > IP
* Bloqueio temporário após exceder o limite
* Persistência em Redis
* Middleware HTTP
* Configuração via `.env`
* Docker Compose
* Testes automatizados

---

## 🏗️ Arquitetura

```text
HTTP Request
      │
      ▼
Middleware
      │
      ▼
RateLimiter (Use Case)
      │
      ▼
RateLimiterRepository
      │
      ▼
RedisRepository
      │
      ▼
Redis
```

### Decisões de Projeto

**Separação de responsabilidades**

* Middleware: captura IP/Token e delega a decisão.
* Use Case: contém toda a regra de negócio.
* Repository: abstrai a persistência.

**Strategy Pattern**

A persistência é acessada através da interface:

```go
type RateLimiterRepository interface {
	IsBlocked(ctx context.Context, key string) (bool, error)
	Increment(ctx context.Context, key string, window time.Duration) (int64, error)
	Block(ctx context.Context, key string, duration time.Duration) error
}
```

A implementação atual utiliza Redis, mas novas estratégias podem ser adicionadas sem alterar a regra de negócio.

---

## ⚙️ Configuração

Arquivo `.env`

```env
PORT=8080

REDIS_ADDR=redis:6379
REDIS_PASSWORD=

RATE_LIMIT_IP=10

BLOCK_DURATION_SECONDS=300

TOKEN_PREMIUM=100
TOKEN_ADMIN=1000
```

### Variáveis

| Variável               | Descrição                   |
| ---------------------- | --------------------------- |
| PORT                   | Porta da aplicação          |
| REDIS_ADDR             | Endereço do Redis           |
| REDIS_PASSWORD         | Senha do Redis              |
| RATE_LIMIT_IP          | Limite padrão por IP        |
| BLOCK_DURATION_SECONDS | Tempo de bloqueio           |
| TOKEN_*                | Limite específico por token |

---

## 🚀 Executando

Subir aplicação e Redis:

```bash
docker compose up --build
```

Aplicação disponível em:

```text
http://localhost:8080
```

---

## 🧪 Testes

Executar todos os testes:

```bash
go test ./...
```

### Cenários cobertos

* Permite requisições abaixo do limite
* Bloqueia requisições acima do limite
* Precedência Token > IP
* Cliente bloqueado continua bloqueado
* Token desconhecido utiliza limite de IP
* IPs possuem contadores independentes

---

## 📌 Exemplos

### Sem Token

```bash
curl http://localhost:8080
```

### Com Token

```bash
curl -H "API_KEY: PREMIUM" http://localhost:8080
```

### PowerShell

```powershell
Invoke-WebRequest `
  -Uri "http://localhost:8080" `
  -Headers @{
      API_KEY = "PREMIUM"
  }
```

---

## 🔒 Comportamento de Bloqueio

Quando o limite configurado é excedido:

**Status HTTP**

```http
429 Too Many Requests
```

**Mensagem**

```text
you have reached the maximum number of requests or actions allowed within a certain time frame
```

O IP ou Token permanece bloqueado durante o período definido por:

```env
BLOCK_DURATION_SECONDS
```

---

## 📁 Estrutura

```text
.
├── cmd/
│   └── server/
│
├── internal/
│   ├── config/
│   ├── gateway/
│   ├── infra/
│   │   └── redis/
│   ├── middleware/
│   └── usecase/
│
├── Dockerfile
├── docker-compose.yml
├── .env
└── README.md
```

---

## Tecnologias

* Go
* Redis
* Docker
* Docker Compose
* go-redis/v9
