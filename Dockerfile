FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /marketplace cmd/api/main.go

FROM scratch
COPY --from=builder /marketplace /marketplace
EXPOSE 8080
ENTRYPOINT ["/marketplace"]
