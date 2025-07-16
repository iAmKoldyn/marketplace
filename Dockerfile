############################################
# Build stage
############################################
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the API server
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api    cmd/api/main.go

# Build the worker
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/worker cmd/worker/main.go

############################################
# Final stage (scratch)
############################################
FROM scratch
# Copy both executables in
COPY --from=builder /bin/api    /api
COPY --from=builder /bin/worker /worker

# Expose API port
EXPOSE 8080

# Default entrypoint is the API; override in docker‑compose for the worker
ENTRYPOINT ["/api"]
