# Usar uma imagem base do Go
FROM golang:1.23-bullseye AS builder

RUN apt-get update && \
	apt-get upgrade -y && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*
# Diretório de trabalho no container
WORKDIR /app

# Copiar arquivos para o container
COPY go.mod go.sum ./
RUN go mod download

# Copiar o código-fonte
COPY . .


# Compilar o binário
RUN go build -o sensor .

# Segunda etapa: criar uma imagem leve para execução
FROM debian:bullseye

RUN apt-get update && \
    apt-get install -y ca-certificates netcat && \
    rm -rf /var/lib/apt/lists/*

# Diretório de trabalho no container
WORKDIR /app

# Copiar o binário da etapa de build
COPY --from=builder /app/sensor .
COPY data /app/data

# Adicionar script de espera
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Expor a porta (se necessário para monitoramento)
EXPOSE 1883

# Comando padrão para rodar o binário
ENTRYPOINT ["/app/entrypoint.sh"]

CMD ["./sensor"]
