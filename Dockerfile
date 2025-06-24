FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o piro-bot .


FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/piro-bot .
RUN chmod 755 ./piro-bot

CMD ["./piro-bot"]