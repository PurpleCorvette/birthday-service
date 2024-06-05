FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /birthday-service cmd/main.go

FROM debian:buster-slim

WORKDIR /root/

COPY --from=builder /birthday-service .

COPY .env .

CMD ["./birthday-service"]