# Usar uma imagem base Alpine
FROM alpine:latest

WORKDIR /app

# Instalar dependências necessárias
RUN apk add --no-cache wget unzip

# Baixar e extrair o binário pré-compilado do PocketBase
RUN wget -O pocketbase.zip https://github.com/pocketbase/pocketbase/releases/download/v0.25.1/pocketbase_0.25.1_linux_amd64.zip && \
    unzip pocketbase.zip && \
    rm pocketbase.zip && \
    chmod +x /app/pocketbase

# Expor a porta padrão do PocketBase (8090)
EXPOSE 8090

# Rodar o PocketBase
CMD ["./pocketbase", "serve", "--http=0.0.0.0:8090"]
