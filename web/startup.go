package web

import (
	"net/http"

	"github.com/autograde/srv/config"
	"github.com/julienschmidt/httprouter"
)

// form value constants used to ensure consistent naming
const (
	fvHomePage     = "homePageURL"
	fvAdmin        = "admin"
	fvClientID     = "clientID"
	fvClientSecret = "clientSecret"
	fvStorageLoc   = "storageLoc"
)

// startup handles requests for the system setup page.
// This is called when database is empty.
func startup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	startupView := struct {
		OptionalHeadline bool
		SysName          string
		SysNameLC        string
		StoragePath      string
		HomePageURL      string
		SaveConfigPath   string
		HomePage         string
		Admin            string
		ClientID         string
		ClientSecret     string
		StorageLoc       string
	}{
		false, config.SysName, config.SysNameLC, config.StdPath, config.SysURL,
		saveConfigPath, fvHomePage, fvAdmin, fvClientID, fvClientSecret, fvStorageLoc,
	}
	execTemplate("startup.html", w, startupView)
}

var saveConfigPath = "/save/config"

// saveConfig handles post requests to save the configuration data provided
// in the system configuration form. This method will redirect if
func saveConfig(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if config.IsSet() {
		logAndRedirect(w, r, front, "Configuration already initialized")
		return
	}
	url := r.FormValue(fvHomePage)
	admin := r.FormValue(fvAdmin)
	clientID := r.FormValue(fvClientID)
	clientSecret := r.FormValue(fvClientSecret)
	path := r.FormValue(fvStorageLoc)
	conf, err := config.NewConfig(url, admin, clientID, clientSecret, path)
	if err != nil {
		// can't continue without proper configuration
		logErrorAndRedirect(w, r, startupPage, err)
	}
	if err := conf.Save(); err != nil {
		// can't continue without proper configuration
		logErrorAndRedirect(w, r, startupPage, err)
	}
	// set global configuration struct; will be accessible through config.Get()
	conf.SetCurrent()
	logAndRedirect(w, r, home, "Configuration initialized")
}
