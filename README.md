# 🚀 Flyer API

API de flyers/panfletos digitais construída com **Golang** e **PocketBase**.

## 🛠️ Tech Stack
- **Backend:** Golang (Go 1.25+)
- **Router:** Chi (go-chi/chi/v5)
- **Database:** SQLite (modernc.org/sqlite)
- **BaaS:** PocketBase (v0.25.1)
- **Containerização:** Docker Multi-stage + Docker Compose

## 📋 Sobre o Projeto
Sistema de API para gerenciamento de flyers/panfletos digitais, permitindo que estabelecimentos criem e distribuam material promocional de forma automatizada. Parte do ecossistema Flyer que inclui apps Flutter e interfaces web.

## 🚀 Como Rodar

### Com Docker
```bash
docker-compose up -d
```

A API Go estará disponível em `http://localhost:8080/api`.
O PocketBase estará disponível em `http://localhost:8090/_/`.

### Sem Docker
```bash
# Iniciar PocketBase primeiro
./pocketbase serve --http=0.0.0.0:8090

# Em outro terminal, iniciar a API Go
cd app/cmd/server && go run main.go
```

## 📁 Estrutura
```
flyer/
├── app/
│   ├── cmd/server/main.go        # Entry point da API
│   └── internal/
│       ├── auth/middleware.go     # Autenticação via API Key
│       ├── database/db.go         # Conexão SQLite + migrations
│       ├── handlers/
│       │   ├── flyers.go          # CRUD flyers
│       │   ├── estabelecimentos.go # CRUD estabelecimentos
│       │   ├── categorias.go      # CRUD categorias
│       │   ├── upload.go          # Upload de imagens
│       │   └── helpers.go         # Utilitários
│       ├── middleware/
│       │   ├── cors.go            # CORS middleware
│       │   └── logging.go         # Logging middleware
│       ├── models/
│       │   ├── flyer.go
│       │   ├── estabelecimento.go
│       │   └── categoria.go
│       └── router/
│           └── router.go          # Configuração de rotas (Chi)
├── Dockerfile                     # Build multi-stage
├── docker-compose.yml             # Orquestração dos containers
├── go.mod
└── go.sum
```

## 🔌 Endpoints da API

### Health Check
```
GET /api/health
```

### Flyers
```
GET    /api/flyers              # Listar todos
GET    /api/flyers/{id}         # Obter um
POST   /api/flyers              # Criar
PUT    /api/flyers/{id}         # Atualizar
DELETE /api/flyers/{id}         # Remover
```

### Estabelecimentos
```
GET    /api/estabelecimentos              # Listar todos
GET    /api/estabelecimentos/{id}         # Obter um
POST   /api/estabelecimentos              # Criar
PUT    /api/estabelecimentos/{id}         # Atualizar
DELETE /api/estabelecimentos/{id}         # Remover
```

### Categorias
```
GET    /api/categorias              # Listar todas
GET    /api/categorias/{id}         # Obter uma
POST   /api/categorias              # Criar
PUT    /api/categorias/{id}         # Atualizar
DELETE /api/categorias/{id}         # Remover
```

### Upload de Imagens
```
POST   /api/upload              # Upload (multipart/form-data, campo: file)
GET    /uploads/{filename}      # Servir arquivo
```

## 🔐 Autenticação
Configure a variável de ambiente `FLYER_API_KEY` para ativar a proteção por API Key.
Envie a chave no header `X-API-Key` ou no query parameter `api_key`.
Sem a variável configurada, todas as requisições são permitidas.

## 🔗 Projetos Relacionados
- [flyer_only_flutter](https://github.com/marcuslinhares/flyer_only_flutter) — App mobile em Flutter
- [flyer-estabelecimento-ui](https://github.com/marcuslinhares/flyer-estabelecimento-ui) — Interface para estabelecimentos
- [flyer-ui-web](https://github.com/marcuslinhares/flyer-ui-web) — Interface web

## 📝 Licença
MIT
