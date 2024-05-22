FROM golang:1.17 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest AS publish
WORKDIR /root/
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]