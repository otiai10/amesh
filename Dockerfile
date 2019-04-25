FROM golang:latest

RUN go get -u github.com/otiai10/amesh

CMD ["amesh"]