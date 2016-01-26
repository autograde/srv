package web

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/golang/glog"
)

func logServerError(w http.ResponseWriter, err error) {
	errCode := rand.Uint32()
	errMsg := fmt.Sprintf("InternalServerError(%d)", errCode)
	http.Error(w, errMsg, http.StatusInternalServerError)
	glog.ErrorDepth(2, errMsg, ": ", err)
}

func logNotFoundError(w http.ResponseWriter, err error) {
	errCode := rand.Uint32()
	errMsg := fmt.Sprintf("NotFound(%d)", errCode)
	http.Error(w, errMsg, http.StatusNotFound)
	glog.ErrorDepth(1, errMsg, ": ", err)
}

func logErrorAndRedirect(w http.ResponseWriter, r *http.Request, redirectTo string, err error) {
	http.Redirect(w, r, redirectTo, http.StatusTemporaryRedirect)
	glog.ErrorDepth(1, err)
}

func logAndRedirect(w http.ResponseWriter, r *http.Request, redirectTo, msg string) {
	http.Redirect(w, r, redirectTo, http.StatusTemporaryRedirect)
	glog.InfoDepth(1, msg)
}
