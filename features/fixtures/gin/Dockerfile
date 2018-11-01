
ARG GO_VERSION
FROM bugsnag_go_test:$GO_VERSION

ARG GIN_VERSION
RUN (cd /app/src/github.com/gin-gonic/gin && git checkout $GIN_VERSION && go install)

COPY . /app/src/test
WORKDIR /app/src/test/scenarios
