ARG GO_VERSION
FROM golang:${GO_VERSION}-alpine

RUN apk update && \
    apk upgrade && \
    apk add git

ENV GOPATH /app

COPY testbuild /app/src/github.com/bugsnag/bugsnag-go
WORKDIR /app/src/github.com/bugsnag/bugsnag-go/v2

# Ensure subsequent steps are re-run if the GO_VERSION variable changes
ARG GO_VERSION

# Get bugsnag dependencies using a conditional call to run go get or go install based on the go version
RUN if [[ $(echo -e "1.11\n$GO_VERSION\n1.16" | sort -V | head -2 | tail -1) == "$GO_VERSION" ]]; then \
        echo "Version is between 1.11 and 1.16, running go get"; \
        go get ./...; \
    else \
        echo "Version is greater than 1.16, running go install"; \
        go install ./...; \
    fi

# Copy test scenarios
COPY ./net_http /app/src/test
WORKDIR /app/src/test

# Create app module - avoid locking bugsnag dep by not checking it in
# Skip on old versions of Go which pre-date modules
RUN if [[ $GO_VERSION != '1.11' && $GO_VERSION != '1.12' ]]; then \
        go mod init && go mod tidy; \
    fi
