FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o backend main.go

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/backend /app/backend
COPY .env /app/.env
COPY models /app/models
COPY handlers /app/handlers
COPY utils /app/utils
EXPOSE 8080
CMD ["/app/backend"]
