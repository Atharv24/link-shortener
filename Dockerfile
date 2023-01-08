FROM golang:latest

ADD . /go/src/link-shortener

WORKDIR /go/src/link-shortener

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["link-shortener"]
