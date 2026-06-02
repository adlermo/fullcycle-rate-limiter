FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o server ./cmd/server

FROM alpine

WORKDIR /root

COPY --from=builder /app/server .

CMD ["./server"]