FROM golang:1.23

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o adres-talk ./

CMD ["./adres-talk $TOKEN"]