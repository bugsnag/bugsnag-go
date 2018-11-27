# Bugsnag-Go Maze-Runner tests

These are feature tests, built on top of [maze-runner](https://github.com/bugsnag/maze-runner) - a Cucumber wrapper with convenience steps for testing Bugsnag notifiers.

In order to run these tests locally you will need a Unix shell, Docker (and docker-compose) and Bundle installed.

You can then run all the tests locally using the `run-maze.sh` script located in this directory from the root of the repository.

```bash
features/run-maze.sh
```

## Running specific features

You can run the maze-tests on a feature-by-feature basis too.

To run only a single feature you can do the following from the root of the repository.

```bash
bundle install #only needs to be done once
# Only run the appversion feature for negroni
GO_VERSION=1.11 NEGRONI_VERSION=v1.0.0 bundle exec bugsnag-maze-runner features/negroni_features/appversion.feature
```

Note that you will have to specify the Go and framework versions for each call.
Also note that when testing revel you'll have to specify both the `REVEL_VERSION` and `REVEL_CMD_VERSION`.
For martini we only support version 1.0, so no `MARTINI_VERSION` variable needs to be set.

Similarly, you can also run the tests for one framework:

```bash
# Only run the negroni features
GO_VERSION=1.11 NEGRONI_VERSION=v1.0.0 bundle exec bugsnag-maze-runner features/negroni_features
```
