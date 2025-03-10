FROM --platform=$BUILDPLATFORM golang:1.24.1 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy

ARG TARGETOS
ARG TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -o t-cleaner .

FROM --platform=$TARGETPLATFORM alpine:latest
USER 1000:1000
WORKDIR /app/
COPY --from=builder /app/t-cleaner .
CMD ["/app/t-cleaner"]
