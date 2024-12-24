# Build Stage
FROM golang:1.20 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o backend .

# Run Stage
FROM alpine:3.18
WORKDIR /root/
COPY --from=builder /app/backend .
EXPOSE 8080
CMD ["./backend"]
