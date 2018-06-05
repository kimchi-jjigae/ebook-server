package main

import (
	"log"
	"net/http"
	"time"

	"github.com/husobee/vestigo"
)

func main() {
	router := vestigo.NewRouter()
	// you can enable trace by setting this to true
	vestigo.AllowTrace = true

	// Setting up router global CORS policy
	// These policy guidelines are overriddable at a per resource level shown below
	router.SetGlobalCors(&vestigo.CorsAccessControl{
		AllowOrigin:      []string{"*", "test.com"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"X-Header", "X-Y-Header"},
		MaxAge:           3600 * time.Second,
		AllowHeaders:     []string{"X-Header", "X-Y-Header"},
	})

	RegisterRoutes(router)

	// Below Applies Local CORS capabilities per Resource (both methods covered)
	// by default this will merge the "GlobalCors" settings with the resource
	// cors settings.  Without specifying the AllowMethods, the router will
	// accept any Request-Methods that have valid handlers associated
	router.SetCors("/welcome", &vestigo.CorsAccessControl{
		AllowMethods: []string{"GET"},                    // only allow cors for this resource on GET calls
		AllowHeaders: []string{"X-Header", "X-Z-Header"}, // Allow this one header for this resource
	})

    port := ":1234"
    log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(port, router))
	//log.Fatal(http.ListenAndServeTLS(":1234", "MyCertificate.crt", "MyKey.key", router))
}
