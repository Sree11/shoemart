#!/usr/bin/env bash

FROM golang
WORKDIR /go/src/esapp
ADD . /go/src/esapp
RUN go get github.com/elastic/go-elasticsearch
RUN go get github.com/gorilla/mux
RUN go get github.com/gorilla/sessions
RUN go get golang.org/x/crypto/bcrypt
RUN go install -v ./...

ENTRYPOINT /go/bin/esapp

EXPOSE 8000