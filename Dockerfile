FROM golang:1.16.3

WORKDIR /go/src/tks-contract
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
