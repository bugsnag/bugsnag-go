TEST?=./...

default: alldeps test

deps:
	go get -v -d ./...

alldeps:
	go get -v -d -t ./...

updatedeps:
	go get -v -d -u ./...

test: alldeps
	#TODO: 2018-09-20 Not testing the 'errors' package as it relies on some very runtime-specific implementation details.
	# The testing of 'errors' needs to be revisited
	go test . ./gin ./martini ./negroni ./sessions ./headers
	@go vet 2>/dev/null ; if [ $$? -eq 3 ]; then \
		go get golang.org/x/tools/cmd/vet; \
	fi
	@go vet $(TEST) ; if [ $$? -eq 1 ]; then \
		echo "go-vet: Issues running go vet ./..."; \
		exit 1; \
	fi

maze:
	bundle install
	bundle exec bugsnag-maze-runner

ci: alldeps test maze

bench:
	go test --bench=.*


.PHONY: bin checkversion ci default deps generate releasebin test testacc testrace updatedeps
