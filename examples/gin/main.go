package main

import (
  "github.com/bugsnag/bugsnag-go"
  "github.com/bugsnag/bugsnag-go/gin"
  "github.com/gin-gonic/gin"
  "net/http"
  "os"
)

type App struct {
  errorHandlerConfig bugsnag.Configuration
}

var app = App{
  bugsnag.Configuration{
    APIKey: "YOUR API KEY",
}}

func main() {

    g := gin.Default()

    g.Use(bugsnaggin.AutoNotify(app.errorHandlerConfig))

    g.GET("/crash", performUnhandledCrash)
    g.GET("/handled", performHandledCrash)

    g.Run(":9001") // listen and serve on 0.0.0.0:9001
}

func performUnhandledCrash(c *gin.Context) {
  c.String(http.StatusOK, "OK")
  var a struct{}
  crash(a)
}

func performHandledCrash(c *gin.Context) {
  _, err := os.Open("some_nonexistent_file.txt")
  if err != nil {
    bugsnag.Notify(err, app.errorHandlerConfig)
  }
  c.String(http.StatusOK, "OK")
}

func crash(a interface{}) string {
  return a.(string)
}
