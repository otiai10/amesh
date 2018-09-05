FROM golang:1.5

MAINTAINER otiai10 <otiai10@gmail.com>

ADD . /go/src/github.com/otiai10/amesh
WORKDIR /go/src/github.com/otiai10/amesh
RUN go get ./...

ENTRYPOINT /go/bin/amesh -d
