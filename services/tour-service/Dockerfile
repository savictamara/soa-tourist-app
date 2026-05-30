FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY services/tour-service/go.mod services/tour-service/go.sum ./
RUN go mod download

COPY services/tour-service/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o tour-service .

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/tour-service ./tour-service

EXPOSE 8085 9091

CMD ["./tour-service"]
