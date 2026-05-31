# ============================================================
# Stage 1: Build Go API
# ============================================================
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/server ./app/cmd/server/

# ============================================================
# Stage 2: Runtime with Go API + PocketBase
# ============================================================
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache wget unzip ca-certificates tzdata

WORKDIR /app

# --- PocketBase ---
RUN wget -q -O pocketbase.zip https://github.com/pocketbase/pocketbase/releases/download/v0.25.1/pocketbase_0.25.1_linux_amd64.zip && \
    unzip pocketbase.zip && \
    rm pocketbase.zip && \
    chmod +x /app/pocketbase

# --- Go API binary ---
COPY --from=builder /app/server /app/server

# Create required directories
RUN mkdir -p /app/data /app/uploads /app/pb_data

# Expose ports: 8080 (Go API) and 8090 (PocketBase)
EXPOSE 8080 8090

# Default environment variables
ENV DB_PATH=/app/data/flyer.db
ENV UPLOAD_DIR=/app/uploads
ENV PORT=8080

# Start both services using a shell script
COPY <<'EOF' /app/start.sh
#!/bin/sh
# Start PocketBase in background
./pocketbase serve --http=0.0.0.0:8090 --dir=/app/pb_data &
# Start Go API
exec ./server
EOF

RUN chmod +x /app/start.sh

CMD ["/app/start.sh"]
