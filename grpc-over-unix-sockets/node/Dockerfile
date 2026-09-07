# ---- Stage 1: build ----
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor/ ./vendor/

COPY proto/ ./proto/
COPY node/main.go ./node/main.go

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -trimpath -ldflags="-s -w" -o /node ./node

# ---- Stage 2: minimal runtime ----
FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /node /node

USER nobody

ENTRYPOINT ["/node"]
