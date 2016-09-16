package bugsnaggin

import (
  "github.com/bugsnag/bugsnag-go"
  "github.com/gin-gonic/gin"
)

// AutoNotify sends any panics to bugsnag, and then re-raises them.
// You should use this after another middleware that
// returns an error page to the client, for example gin.Recovery().
// The arguments can be any RawData to pass to Bugsnag, most usually
// you'll pass a bugsnag.Configuration object.
func AutoNotify(rawData ...interface{}) gin.HandlerFunc {
  return func(c *gin.Context) {
    r := c.Request

    // create a notifier that has the current request bound to it
    notifier := bugsnag.New(append(rawData, r)...)
    defer notifier.AutoNotify(r)
    c.Next()
  }
}
