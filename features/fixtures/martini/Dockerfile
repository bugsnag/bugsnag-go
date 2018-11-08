ARG GO_VERSION
FROM golang:${GO_VERSION}-alpine as notifier_builder

RUN apk update && \
    apk upgrade && \
    apk add git

ENV GOPATH /app

COPY testbuild /app/src/github.com/bugsnag/bugsnag-go
WORKDIR /app/src/github.com/bugsnag/bugsnag-go

RUN go get -v -d ./...

# Copy test scenarios
COPY ./martini /app/src/test
WORKDIR /app/src/test
