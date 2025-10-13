FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go modules files
COPY go.mod go.sum ./

# Download dependencies first (better for Docker layer caching)
RUN go mod download

# Copy source code
COPY . .

# Build the application with static linking
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api cmd/api/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates mysql-client netcat-openbsd
WORKDIR /app

COPY --from=builder /app/api /app/api
COPY migrations/ ./migrations/
COPY entrypoint.sh ./entrypoint.sh

RUN chmod +x entrypoint.sh && chmod +x api

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]
CMD ["./api"]