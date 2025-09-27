FROM golang:1.24.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build static binary (no glibc dependency)
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# Final stage: small image
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/api .

EXPOSE 8080

CMD ["./api"]
