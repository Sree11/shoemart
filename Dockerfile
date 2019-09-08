#!/usr/bin/env bash

FROM golang:1.12.9
RUN mkdir /searchapp
ADD . /searchapp/
WORKDIR /searchapp
COPY . /searchapp

RUN go build -o search .

CMD [ ./searchapp/search" ]