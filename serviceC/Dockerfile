# Dockerfile for serviceA
FROM golang:1.23.4

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o service_c .

EXPOSE 3000

CMD ["./service_c"]