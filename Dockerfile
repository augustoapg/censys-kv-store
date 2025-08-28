# --- Build stage ---
FROM golang:1.25-alpine AS builder

WORKDIR /app

# copy the source code and download dependencies
COPY . .
RUN go mod download

# create a statically linked binary (no external dependencies)
RUN CGO_ENABLED=0 GOOS=linux go build -o kv-store .

# --- Runtime stage ---
FROM alpine:latest

WORKDIR /app

# copy just the binary from the build stage
COPY --from=builder /app/kv-store .

# create and use non-root user and group to avoid running as root
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup
RUN chown appuser:appgroup kv-store

USER appuser
EXPOSE 8080

CMD ["./kv-store"]
