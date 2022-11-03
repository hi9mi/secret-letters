FROM golang:1.19

ENV GIN_MODE=release 

WORKDIR /usr/src/secret-letters

COPY go.* ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/secret-letters ./...

CMD ["secret-letters"]