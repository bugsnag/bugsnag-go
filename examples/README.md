# Examples of working with bugsnag-go

In this directory you can find example applications of the frameworks we support, and other examples of common use cases.

The examples that expose a HTTP port will all listen on 9001.

## Use cases and frameworks

* [Capturing panics within goroutines](using-goroutines). Goroutines require special care to avoid crashing the app entirely or cleaning up before an error report can be sent.
  This is an example of a panic within a goroutine which is sent to Bugsnag.
* [Using net/http](http) (web server using the standard library)
* [Using Gin](gin) (web framework)
* [Using Negroni](negroni) (web framework)
* [Using Martini](martini) (web framework)
* [Using Revel](revelapp) (web framework)
