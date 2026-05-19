# 🚀 Flyer API

API de flyers/panfletos digitais construída com **Golang** e **PocketBase**.

## 🛠️ Tech Stack
- **Backend:** Golang
- **Database/BaaS:** PocketBase
- **Containerização:** Docker + Docker Compose

## 📋 Sobre o Projeto
Sistema de API para gerenciamento de flyers/panfletos digitais, permitindo que estabelecimentos criem e distribuam material promocional de forma automatizada. Parte do ecossistema Flyer que inclui apps Flutter e interfaces web.

## 🚀 Como Rodar

### Com Docker
```bash
docker-compose up -d
```

### Sem Docker
```bash
go run main.go
```

## 📁 Estrutura
```
flyer/
├── main.go              # Entry point da API
├── Dockerfile           # Build da imagem
└── docker-compose.yml   # Orquestração dos containers
```

## 🔗 Projetos Relacionados
- [flyer_only_flutter](https://github.com/marcuslinhares/flyer_only_flutter) — App mobile em Flutter
- [flyer-estabelecimento-ui](https://github.com/marcuslinhares/flyer-estabelecimento-ui) — Interface para estabelecimentos
- [flyer-ui-web](https://github.com/marcuslinhares/flyer-ui-web) — Interface web

## 📝 Licença
MIT
