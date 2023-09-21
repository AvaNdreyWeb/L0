FROM golang:latest

WORKDIR /app

COPY go.mod go.sum main.go ./
RUN go mod download

RUN go build -o main ./

EXPOSE 8081

CMD ["./main"]