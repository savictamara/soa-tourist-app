FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY services/purchase-service/go.mod services/purchase-service/go.sum ./
RUN go mod download

COPY services/purchase-service/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o purchase-service .

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/purchase-service ./purchase-service

EXPOSE 8087 9093

CMD ["./purchase-service"]
