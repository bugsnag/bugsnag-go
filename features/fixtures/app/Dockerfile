ARG GO_VERSION
FROM golang:${GO_VERSION}-alpine

# Only applicable to development environments running Zscaler
# These Docker commands are ugly but the approach being used across SB
COPY .zscaler-root-ca.pem* /tmp/zscaler-root-ca.crt
RUN if [ -f /tmp/zscaler-root-ca.crt ]; then \
        echo "Zscaler root CA certificate found, adding to trusted certificates"; \
        set -eux; \
        if [ -f /etc/ssl/cert.pem ]; then \
            cat /tmp/zscaler-root-ca.crt >> /etc/ssl/cert.pem; \
        elif [ -f /etc/ssl/certs/ca-certificates.crt ]; then \
            cat /tmp/zscaler-root-ca.crt >> /etc/ssl/certs/ca-certificates.crt; \
        else \
            mkdir -p /etc/ssl/certs; \
            cp /tmp/zscaler-root-ca.crt /etc/ssl/certs/ca-certificates.crt; \
            ln -sf /etc/ssl/certs/ca-certificates.crt /etc/ssl/cert.pem; \
        fi; \
        apk add --no-cache ca-certificates; \
        mkdir -p /usr/local/share/ca-certificates; \
        cp /tmp/zscaler-root-ca.crt /usr/local/share/ca-certificates/zscaler-root-ca.crt; \
        chmod 644 /usr/local/share/ca-certificates/zscaler-root-ca.crt; \
        update-ca-certificates; \
        rm -f /tmp/zscaler-root-ca.crt \
    else \
        echo "No Zscaler root CA certificate found, skipping addition to trusted certificates"; \
    fi;

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