# Build stage
FROM golang:1.21-alpine3.19 AS build

WORKDIR /app/build

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o codespace-backend-api ./cmd/main.go

# Publish stage
FROM alpine:3.19 AS publish

ARG DEFAULT_PORT=80

ENV PORT=${DEFAULT_PORT}
ENV HTTP_PORT=${DEFAULT_PORT}

ENV HTTP_ADDRESS=0.0.0.0
ENV TZ=America/Sao_Paulo

EXPOSE ${DEFAULT_PORT}

WORKDIR /api

COPY --from=build /app/build/infra/deploy/entrypoint.sh /api
COPY --from=build /app/build/codespace-backend-api /api

RUN chmod +x ./codespace-backend-api
RUN chmod +x ./entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]