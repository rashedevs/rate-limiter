FROM golang:1.24.3-alpine

WORKDIR /app

COPY . .

RUN go build -o server .

EXPOSE 8058

CMD ["./server"]
