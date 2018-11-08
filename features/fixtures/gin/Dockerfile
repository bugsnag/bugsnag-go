ARG GO_VERSION
FROM golang:${GO_VERSION}-alpine as notifier_builder

RUN apk update && \
    apk upgrade && \
    apk add git

ENV GOPATH /app

COPY testbuild /app/src/github.com/bugsnag/bugsnag-go
WORKDIR /app/src/github.com/bugsnag/bugsnag-go

RUN go get -v -d ./...

FROM notifier_builder

# Switch to correct version of gin
ARG GIN_VERSION
RUN test -n "$GIN_VERSION"

RUN (cd /app/src/github.com/gin-gonic/gin && git checkout $GIN_VERSION)
RUN go get -v -d ./...
RUN (cd /app/src/github.com/gin-gonic/gin && go install)


# Copy test scenarios
COPY ./gin /app/src/test
WORKDIR /app/src/test