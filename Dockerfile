# Stage 1: Build
# Uses the full Go image to compile the binary
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy dependency files first (cached as a layer — only re-runs if go.mod/go.sum change)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o server .

# Stage 2: Run
# Tiny image — just the binary, nothing else
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
