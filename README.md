# cors

Negroni middleware/handler to enable CORS support.


## Usage

~~~go
import (
       "github.com/codegangsta/negroni"
       "github.com/hariharan-uno/cors"
)

func main() {
     n := negroni.Classic()

     // CORS for https://*.foo.com origins, allowing:
     // - GET and POST methods
     // - Origin header 
     opts := cors.Options{
     	  AllowOrigins:    []string{"https://*.foo.com"},
	  AllowMethods:    []string{"GET", "POST"},
	  AllowHeaders:	   []string{"Origin"},
     }

     n.Use(negroni.HandlerFunc(opts.Allow))

     mux := http.NewServeMux()
     // map your routes

     n.UseHandler(mux)

     n.Run(":3000")
}
~~~

## Authors

* [Burcu Dogan](http://github.com/rakyll)
* [Hari haran](http://github.com/hariharan-uno)