FROM golang:1.25 AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/server .

RUN mkdir -p /data

EXPOSE 8080

ENV LISTEN_ADDR=0.0.0.0:8080
ENV DATA_PATH=/data
ENV TOKEN=123456

CMD ["./server"]