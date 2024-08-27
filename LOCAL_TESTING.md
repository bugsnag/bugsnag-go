
## Unit tests
* Install old golang version (do not install just 1.11 - it's not compatible with running newer modules): 

```
ASDF_GOLANG_OVERWRITE_ARCH=amd64 asdf install golang 1.11.13
```

* If you see error below use `CGO_ENABLED=0`.

```
# crypto/x509
malformed DWARF TagVariable entry
```

## Local testing with maze runner

* Maze runner tests require
  * Specyfing `GO_VERSION` env variable to set a golang version for docker container.
  * Ruby 2.7.
  * Running docker.

* Commands to run tests

```
	bundle install
	bundle exec bugsnag-maze-runner
  bundle exec bugsnag-maze-runner -c features/<chosen_feature>
```