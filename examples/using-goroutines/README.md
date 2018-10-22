# Managing panics from goroutines

This package contains an example for how to manage panics in a separate goroutine with Bugsnag.

## Run the example

1. Change the API key in `main.go` to a project you've created in [Bugsnag](https://app.bugsnag.com).
1. Inside `bugsnag-go/examples/using-goroutines` do:
    ```bash
    go get
    go run main.go
    ```
1. The application will run for a split second, starting a new goroutine, which panics.
1. You should now see events this panic in your [Bugsnag dashboard](https://app.bugsnag.com).
