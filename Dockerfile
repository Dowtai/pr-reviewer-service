FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o pr-reviewer-service ./main
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/pr-reviewer-service .
EXPOSE 8080
CMD ["./pr-reviewer-service"]
