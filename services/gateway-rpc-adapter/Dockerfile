FROM golang:1.24-alpine AS build

WORKDIR /app
COPY services/gateway-rpc-adapter/go.mod ./
COPY services/gateway-rpc-adapter/go.sum* ./
RUN go mod download

COPY services/gateway-rpc-adapter/ ./
RUN go build -o gateway-rpc-adapter .

FROM alpine:3.20

WORKDIR /app
COPY --from=build /app/gateway-rpc-adapter .
EXPOSE 8090

ENTRYPOINT ["./gateway-rpc-adapter"]
