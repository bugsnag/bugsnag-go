ARG GO_VERSION
FROM golang:${GO_VERSION:-1.11}

ENV GOPATH /app

COPY testbuild /app/src/github.com/bugsnag/bugsnag-go
WORKDIR /app/src/github.com/bugsnag/bugsnag-go

RUN go get -v -d ./...

# Copy test scenarios
COPY ./martini /app/src/test
WORKDIR /app/src/test
