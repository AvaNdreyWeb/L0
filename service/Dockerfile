FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg

RUN go build -o main ./cmd/service

EXPOSE 8080

CMD ["./main"]