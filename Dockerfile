# -------- Build Stage --------
  FROM golang:1.24-bullseye as builder
  WORKDIR /app
  RUN apt-get update && apt-get install -y gcc libc6-dev
  COPY go.mod go.sum ./
  RUN go mod download
  COPY . .
  ENV CGO_ENABLED=1
  RUN go build -o proxy-gateway ./cmd/server/main.go
  
  # -------- Runtime Stage --------
  FROM debian:bullseye-slim
  WORKDIR /app
  RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
  COPY --from=builder /app/proxy-gateway .
  COPY --from=builder /app/static ./static
  COPY --from=builder /app/templates ./templates
  EXPOSE 8080
  ENV PORT=8080
  
  CMD ["./proxy-gateway"]
  