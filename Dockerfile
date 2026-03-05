# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Git required for go mod download when fetching modules from VCS (e.g. x-proto)
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o /xword ./cmd/xword

# Runtime stage
FROM gcr.io/distroless/base-debian12
COPY --from=builder /xword /xword
EXPOSE 8080
ENV GRPC_PORT=8080
USER 65534:65534
ENTRYPOINT ["/xword"]
