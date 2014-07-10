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

     opts := cors.Options{
     	  AllowAllOrigins: true,
	  AllowMethods:    []string{"GET", "POST"},
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