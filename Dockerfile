FROM golang:latest AS build-env

ADD src $GOPATH/src

WORKDIR $GOPATH/src/main

RUN go build

RUN mkdir /dist

RUN cp $GOPATH/src/main/comictracker /dist/comictracker

CMD ["/dist/comictracker"]
