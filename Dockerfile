# Build stage
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o terminal-portfolio .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/terminal-portfolio .

# Create SSH directory for mounting keys
RUN mkdir -p .ssh

EXPOSE 2222

CMD ["./terminal-portfolio"]