ARG GO_VERSION
FROM golang:$GO_VERSION as notifier_builder

WORKDIR /app

ENV GOPATH /app

COPY testbuild /app/src/github.com/bugsnag/bugsnag-go
WORKDIR /app/src/github.com/bugsnag/bugsnag-go

RUN go get -v -d ./...

FROM notifier_builder

# Switch to correct verison of gin
ARG GIN_VERSION
RUN (cd /app/src/github.com/gin-gonic/gin && git checkout $GIN_VERSION && go install)

# Copy test scenarios
COPY ./gin/scenarios /app/src/test
WORKDIR /app/src/test