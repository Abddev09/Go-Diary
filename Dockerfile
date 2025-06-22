# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o blog-server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Static fayllar uchun papkalar
RUN mkdir -p static uploads

COPY --from=builder /app/blog-server .
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./blog-server"]
