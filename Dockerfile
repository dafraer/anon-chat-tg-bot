FROM golang:1.23

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o AnonBot .

CMD ["sh", "-c", "./AnonBot $TOKEN $DB_URI $ADMIN_ID"]
