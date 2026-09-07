FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/server /server

USER nobody

EXPOSE 8080
ENTRYPOINT ["/server"]
