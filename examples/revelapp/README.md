# Example Revel application

This package contains an example Revel application, with Bugsnag configured.

The key files for integrating Bugsnag are:

1. `app/init.go` - Sets up the Bugsnag filter.
1. `app/controllers/app.go` - notifies about an handled error and an unhandled panic
1. `conf/app.conf` - configures Bugsnag, in particular the API key

## Run the example

1. Change the API key in `app.conf` to a project you've created in [Bugsnag](https://app.bugsnag.com).
1. Inside `bugsnag-go/examples/revelapp` do:
    ```bash
    revel run
    ```
1. The application is now running. You can now visit
    ```
    http://localhost:9001/unhandled - to trigger an unhandled panic
    http://localhost:9001/handled   - to trigger a handled error
    ```
1. You should now see events for these exceptions in your [Bugsnag dashboard](https://app.bugsnag.com).
