FROM alpine:latest as build

RUN apk update
RUN apk add \
  git \
  go

RUN go install -v github.com/otiai10/amesh@latest

FROM alpine:latest AS exec
RUN apk add tzdata
COPY --from=build /root/go/bin/amesh /bin
CMD ["amesh"]
