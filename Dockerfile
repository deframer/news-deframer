FROM golang:1.25-alpine

WORKDIR /app

COPY main.go .

# Initialize a module (required for modern Go) and build
RUN go mod init deframer && go build -o app main.go

CMD ["./app"]
