FROM golang:1.19-alpine AS builder

WORKDIR /app

# COPY go.mod go.sum ./
COPY . ./


RUN go mod download

COPY ./otp_service .

RUN go build -o /app/otp_service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/otp_service /app/otp_service

RUN chmod +x /app/otp_service

EXPOSE 50052

CMD ["./otp_service"]
