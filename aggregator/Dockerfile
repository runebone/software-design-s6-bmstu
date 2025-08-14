FROM golang:1.23.1-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum .

RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

ENV PORT 8080

EXPOSE $PORT

CMD ["./main"]
