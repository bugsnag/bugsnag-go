package main

import (
  "github.com/bugsnag/bugsnag-go"
  "github.com/bugsnag/bugsnag-go/gin"
  "github.com/gin-gonic/gin"
  "net/http"
)

func main() {

    g := gin.Default()

    g.Use(bugsnaggin.AutoNotify(bugsnag.Configuration{
      APIKey: "066f5ad3590596f9aa8d601ea89af845"
    }))

    g.GET("/", func(c *gin.Context) {
        performGet(c)
    })

    g.Run(":9001") // listen and serve on 0.0.0.0:9001
}

func performGet(c *gin.Context) {
  c.String(http.StatusOK, "OK")
  var a struct{}
  crash(a)
}

func crash(a interface{}) string {
  return a.(string)
}
