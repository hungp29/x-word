# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o /xword ./cmd/xword

# Runtime stage
FROM scratch
COPY --from=builder /xword /xword
EXPOSE 8080
ENV HTTP_PORT=8080
USER 65534:65534
ENTRYPOINT ["/xword"]
