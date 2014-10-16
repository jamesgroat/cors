# cors  
[![wercker status](https://app.wercker.com/status/bbd90708c634abd65747ca5578de4a90/s "wercker status")](https://app.wercker.com/project/bykey/bbd90708c634abd65747ca5578de4a90)

[Negroni](https://github.com/codegangsta/negroni) middleware/handler to enable CORS support.

## Usage

~~~go
package main

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/hariharan-uno/cors"
)

func main() {
	n := negroni.Classic()

	// CORS for https://*.foo.com origins, allowing:
	// - GET and POST methods
	// - Origin header
	options := cors.Options{
		AllowOrigins: []string{"https://*.foo.com"},
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Origin"},
	}

	n.Use(negroni.HandlerFunc(options.Allow))

	mux := http.NewServeMux()
	// map your routes

	n.UseHandler(mux)

	n.Run(":3000")
}
~~~

###Documentation

[Godoc](https://godoc.org/github.com/hariharan-uno/cors)

## Authors

* [Burcu Dogan](http://github.com/rakyll)
* [Hari haran](http://github.com/hariharan-uno)
