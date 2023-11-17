FROM golang:1.21 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go test -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -C cmd/user-api/ -o /app/main

CMD ["/app/main"]