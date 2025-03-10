FROM golang:1.24.1 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o t-cleaner

FROM alpine:latest
USER 1000:1000
WORKDIR /app/
COPY --from=builder /app/t-cleaner .
CMD ["./t-cleaner"]
