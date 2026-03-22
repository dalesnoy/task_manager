FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/frontend ./frontend

EXPOSE 8080

CMD ["./main"]
