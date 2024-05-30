FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /auth_service cmd/main.go

FROM alpine:latest

COPY --from=builder /auth_service /auth_service

EXPOSE 50051

CMD ["/auth_service"]
