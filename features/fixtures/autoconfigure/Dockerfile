ARG GO_VERSION
FROM golang:${GO_VERSION}-alpine

RUN apk update && apk upgrade && apk add git bash

ENV GOPATH /app

COPY testbuild /app/src/github.com/bugsnag/bugsnag-go
WORKDIR /app/src/github.com/bugsnag/bugsnag-go/v2

# Get bugsnag dependencies
RUN go get ./...

# Copy test scenarios
COPY ./autoconfigure /app/src/test
WORKDIR /app/src/test

# Ensure subsequent steps are re-run if the GO_VERSION variable changes
ARG GO_VERSION
# Create app module - avoid locking bugsnag dep by not checking it in
# Skip on old versions of Go which pre-date modules
RUN if [[ $GO_VERSION != '1.11' && $GO_VERSION != '1.12' ]]; then \
        go mod init && go mod tidy; \
    fi

RUN chmod +x run.sh
CMD ["/app/src/test/run.sh"]
