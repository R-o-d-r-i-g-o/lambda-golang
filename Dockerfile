FROM golang:1.21-alpine3.19 AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest AS publish
WORKDIR /root/
COPY --from=build /app/app .
RUN chmod +x app
EXPOSE 8080

CMD ["./app"]