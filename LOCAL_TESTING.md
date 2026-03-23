
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

## Zscaler environments

If running e2e tests in an environment with Zscaler, first run copy the root certificate into place using:
```
ruby scripts/zscaler_setup.rb
```
The script excepts the certificate to be located at `./.zscaler-root-ca.pem`.

## Local testing with Maze Runner

* Maze Runner tests require
  * Specifying `GO_VERSION` env variable to set a golang version for docker container.
  * Ruby 2.7.
  * Running docker.

* Commands to run tests

```
bundle install
bundle exec maze-runner
bundle exec maze-runner -c features/<chosen_feature>
```