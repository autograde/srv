package web

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// profile handles requests for the profile page of a user.
func profile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("user"))
}

// index handles requests for the system home page.
func index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "greetings, %s!\n", "Autograder")
}

// home handles requests for the home page of a user.
func hom(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "greetings, %s!\n", "Hein Meling")
}
