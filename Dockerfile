FROM golang:1.22.2

WORKDIR /library

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /library/cmd/api

RUN go build -o /library/build/lib

CMD ["/library/build/lib"]