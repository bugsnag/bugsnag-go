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

# Switch to correct version of revel
ARG REVEL_VERSION
RUN test -n "$REVEL_VERSION"

ARG REVEL_CMD_VERSION
RUN test -n "$REVEL_CMD_VERSION"

RUN (cd /app/src/github.com/revel/revel && git checkout $REVEL_VERSION)
RUN (cd /app/src/github.com/revel/revel && go get -v -d ./...)
RUN (cd /app/src/github.com/revel/revel && go install)

RUN go get github.com/revel/cmd/revel
RUN (cd /app/src/github.com/revel/cmd/revel && git checkout $REVEL_CMD_VERSION)
RUN (cd /app/src/github.com/revel/cmd/revel && go get -v -d ./...)
RUN (cd /app/src/github.com/revel/cmd/revel && go install)

# Copy test scenarios
COPY ./revel /app/src/test
WORKDIR /app/src