FROM golang:1.22.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o payment_service ./cmd

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/payment_service .

COPY .env ./

CMD ["./payment_service"]
