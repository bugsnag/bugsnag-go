<div align="center">
  <a href="https://www.bugsnag.com/platforms/android">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://assets.smartbear.com/m/3dab7e6cf880aa2b/original/BugSnag-Repository-Header-Dark.svg">
      <img alt="SmartBear BugSnag logo" src="https://assets.smartbear.com/m/3945e02cdc983893/original/BugSnag-Repository-Header-Light.svg">
    </picture>
  </a>
  <h1>Error monitoring and reporting for Go</h1>
</div>

[![Documentation](https://img.shields.io/badge/documentation-latest-blue.svg)](https://docs.bugsnag.com/performance/go/)
[![Go Reference](https://pkg.go.dev/badge/github.com/bugsnag/bugsnag-go.svg)](https://pkg.go.dev/github.com/bugsnag/bugsnag-go)
[![Build status](https://github.com/bugsnag/bugsnag-go/actions/workflows/test-package.yml/badge.svg?branch=master)](https://buildkite.com/bugsnag/bugsnag-go)

Automatically detect crashes and report errors in your Go apps. Get alerts about errors and panics in real-time, including detailed error reports with diagnostic information. Understand and resolve issues as fast as possible.

Learn more about BugSnag's [Go error monitoring and error reporting](https://www.bugsnag.com/platforms/go-lang-error-reporting/) solution.

## Features

* Automatically report unhandled errors and panics
* Report handled errors
* Attach user information to determine how many people are affected by a crash
* Send customized diagnostic data

## Getting Started

1. [Create a BugSnag account](https://bugsnag.com)
2. Complete the instructions in the integration guide for your framework:
    * [Gin](https://docs.bugsnag.com/platforms/go/gin/)
    * [Negroni](https://docs.bugsnag.com/platforms/go/negroni/)
    * [net/http](https://docs.bugsnag.com/platforms/go/net-http/)
    * [Revel](https://docs.bugsnag.com/platforms/go/revel/)
    * [Other Go apps](https://docs.bugsnag.com/platforms/go/other/)
3. Relax!

## Support

* [Search open and closed issues](https://github.com/bugsnag/bugsnag-go/issues?utf8=✓&q=is%3Aissue) for similar problems
* [Report a bug or request a feature](https://github.com/bugsnag/bugsnag-go/issues/new)

## Contributing

All contributors are welcome! For information on how to build, test and release `bugsnag-go`, see our [contributing guide](CONTRIBUTING.md).


## License

The BugSnag error reporter for Go is free software released under the MIT License. See [LICENSE.txt](LICENSE.txt) for details.
