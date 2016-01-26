package web

//go:generate go-bindata -pkg $GOPACKAGE -o assets_gen.go assets/...

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	front       = "/"
	home        = "/home"
	oauth       = "/oauth"
	signin      = "/login"
	signout     = "/logout"
	startupPage = "/startup"
)

// Start starts the web server.
func Start() {
	router := httprouter.New()
	// Login and OAuth handlers
	// router.Handler("GET", signin, loginHandler())
	// router.HandlerFunc("GET", oauth, oauthHandler)
	// router.HandlerFunc("GET", signout, logoutHandler)

	router.GET(front, index)
	router.GET(home, hom) //TODO renaem
	router.GET("/profile/:user", profile)
	router.GET(startupPage, startup)
	router.POST(saveConfigPath, saveConfig)
	router.Handler("GET", "/new", appHandler())

	log.Fatal(http.ListenAndServe(":80", router))
}
