FROM golang:1.21.5 AS builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./main-app ./cmd/app/

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main-app ./main-app
COPY migrations ./migrations
COPY .env ./.env

EXPOSE 8080
CMD "./main-app"