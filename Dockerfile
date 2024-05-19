FROM golang:1.20-alpine as dev


WORKDIR /root
COPY ./build/config .kube/config

WORKDIR /go/src/app
ENV GO111MODULE=on
COPY . /go/src/app

RUN  go mod download
RUN go build -o main

