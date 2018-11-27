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

# Switch to correct version of negroni
ARG NEGRONI_VERSION
RUN test -n "$NEGRONI_VERSION"

RUN (cd /app/src/github.com/urfave/negroni && git checkout $NEGRONI_VERSION)
RUN go get -v -d ./...
RUN (cd /app/src/github.com/urfave/negroni && go install)


# Copy test scenarios
COPY ./negroni /app/src/test
WORKDIR /app/src/test