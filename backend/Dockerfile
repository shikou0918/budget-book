FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go modules files
COPY go.mod go.sum ./

# Download dependencies first (better for Docker layer caching)
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main cmd/api/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]