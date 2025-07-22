FROM golang:1.24.5-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o order-service ./main.go

EXPOSE 8080
CMD ["./order-service"]
