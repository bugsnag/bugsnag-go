# Example `net/http` application

This package contains an example `net/http` application, with Bugsnag configured.

## Run the example

1. Change the API key in `main.go` to a project you've created in [Bugsnag](https://app.bugsnag.com).
1. Inside `bugsnag-go/examples/http` do:
    ```bash
    go get
    go run main.go
    ```
1. The application is now running. You can now visit
    ```
    http://localhost:9001/unhandled - to trigger an unhandled panic
    http://localhost:9001/handled   - to trigger a handled error
    ```
1. You should now see events for these exceptions in your [Bugsnag dashboard](https://app.bugsnag.com).
