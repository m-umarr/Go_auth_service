FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /otp_service cmd/main.go

FROM alpine:latest

COPY --from=builder /otp_service /otp_service

EXPOSE 50052

CMD ["/otp_service"]
