package main

import (
	"log"
	"net/http"
	"time"

    "github.com/BurntSushi/toml"
	"github.com/husobee/vestigo"
)

type Config struct {
    Password string
    SearchDirs []string
    StorageDir string
    TempDir string
    AllowOriginFrom []string
    Port string
    Certificate string
    Key string
}

func main() {
    var config Config
    if _, err := toml.DecodeFile("config.toml", &config); err != nil {
        log.Fatalf("Misconfigured TOML config file!, %s", err)
    }

    ebooksRouter := newEbooksRouter(&config)
	router := vestigo.NewRouter()
	// you can enable trace by setting this to true
	vestigo.AllowTrace = true

	// Setting up router global CORS policy
	// These policy guidelines are overriddable at a per resource level shown below
	router.SetGlobalCors(&vestigo.CorsAccessControl{
		AllowOrigin:      config.AllowOriginFrom
		AllowCredentials: true,
		MaxAge:           3600 * time.Second,
		AllowHeaders:     []string{"x-password", "content-type"},
	})

	ebooksRouter.RegisterRoutes(router)

    log.Printf("Listening on port %s", config.Port)
	//log.Fatal(http.ListenAndServe(config.Port, router))
	log.Fatal(http.ListenAndServeTLS(config.Port, config.Certificate, config.Key, router))
}
