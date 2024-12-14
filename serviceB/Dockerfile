# Dockerfile for serviceA
FROM golang:1.23.4

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o service_b .

EXPOSE 8081

CMD ["./service_b"]