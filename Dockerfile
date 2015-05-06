FROM golang:1.4.2-cross
MAINTAINER peter.edge@gmail.com

RUN mkdir -p /go/src/github.com/peter-edge/go-ledge
ADD . /go/src/github.com/peter-edge/go-ledge/
WORKDIR /go/src/github.com/peter-edge/go-ledge
