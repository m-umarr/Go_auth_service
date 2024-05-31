FROM golang:1.19-alpine AS builder

WORKDIR /app

# COPY go.mod go.sum ./
COPY . ./
RUN go mod download

COPY ./auth_service .

RUN go build -o /app/auth_service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/auth_service /app/auth_service

RUN chmod +x /app/auth_service

EXPOSE 50051

CMD ["./auth_service"]
