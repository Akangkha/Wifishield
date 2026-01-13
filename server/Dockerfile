# ---------- Build stage ----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY server ./server
COPY agent/proto ./agent/proto


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o netshield-server ./server/cmd/shieldserver

# ---------- Runtime stage ----------
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/netshield-server .

EXPOSE 8082 50051

CMD ["./netshield-server"]
