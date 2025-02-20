FROM golang:1.23

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o bot .

CMD ["sh", "-c", "./bot $TOKEN $DB_URI $ADMIN_ID"]
