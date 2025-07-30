FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o subscription-service ./cmd

FROM alpine:latest
RUN apk add --no-cache ca-certificates bash

WORKDIR /root
COPY --from=builder /app/subscription-service .
COPY migrations ./migrations
COPY deploy/wait-for-it.sh .
RUN chmod +x wait-for-it.sh

EXPOSE 8080

CMD ["./wait-for-it.sh", "db:5432", "--", "./subscription-service"]
