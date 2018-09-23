FROM golang:onbuild

RUN mkdir -p /go/src/shopfiyStoreAPI
WORKDIR /go/src/shopfiyStoreAPI

ADD . /go/src/shopfiyStoreAPI

RUN go get -v
