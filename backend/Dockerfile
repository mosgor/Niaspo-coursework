FROM golang:latest

WORKDIR /backend

COPY . .

EXPOSE 8082

ENV CONFIG_PATH=./config/prod.yaml

RUN go mod tidy

ENTRYPOINT go run cmd/main.go