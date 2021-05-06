FROM golang:1.16.3-alpine

WORKDIR /go/src/app
COPY . .

RUN go install