ARG GO_VERSION
FROM golang:${GO_VERSION}-alpine

RUN apk update && apk upgrade && apk add git bash build-base

ENV GOPATH /app
ENV GO111MODULE="on"

COPY features /app/src/features
COPY v2 /app/src/github.com/bugsnag/bugsnag-go/v2
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

WORKDIR /app/src/features/fixtures/app

# Create app module - avoid locking bugsnag dep by not checking it in
# Skip on old versions of Go which pre-date modules
RUN     go mod init && go mod tidy && \
        echo "replace github.com/bugsnag/bugsnag-go/v2 => /app/src/github.com/bugsnag/bugsnag-go/v2" >> go.mod && \
        go mod tidy

RUN chmod +x run.sh
CMD ["/app/src/features/fixtures/app/run.sh"]