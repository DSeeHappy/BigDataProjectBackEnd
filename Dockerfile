FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o backend main.go

EXPOSE 8080

CMD ["/app/backend"]
