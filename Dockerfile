# ── Stage 1: Build ──
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags "-s -w" -o /bin/perp

# ── Stage 2: Runtime ──
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /bin/perp .

# Web 管理面板端口
EXPOSE 9090
# HTTP API 端口
EXPOSE 8080

ENTRYPOINT ["./perp"]
