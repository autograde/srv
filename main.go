package main

import (
	"flag"
	"log"

	"github.com/autograde/srv/config"
	"github.com/autograde/srv/web"
	"github.com/golang/glog"
)

var (
	admin        = flag.String("admin", "", "Admin must be a valid GitHub username")
	url          = flag.String("url", "", "Homepage URL for "+config.SysName)
	clientID     = flag.String("id", "", "Client ID for OAuth with Github")
	clientSecret = flag.String("secret", "", "Client Secret for OAuth with Github")
	path         = flag.String("path", config.StdPath, "Path for data storage for "+config.SysName)
	help         = flag.Bool("help", false, "Helpful instructions")
)

func main() {
	flag.Parse()

	// set log print appearance
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	if *help {
		// print instructions and command usage
		config.PrintInstructions()
		flag.Usage()
		return
	}

	if *path != "" {
		initConfig()
	}
	glog.Info("Starting webserver")
	web.Start()
}

func initConfig() {
	// load configuration data from the provided base path
	conf, err := config.Load(*path)
	if err != nil {
		glog.Errorln(err)
		// can't load config file; check if other command line arguments provided
		if *admin == "" || *url == "" || *clientID == "" || *clientSecret == "" {
			// no config provided; let user enter through startup page
			return
		}
		// create config based on command line arguments
		conf, err = config.NewConfig(*url, *admin, *clientID, *clientSecret, *path)
		if err != nil {
			// can't continue without proper configuration
			glog.Fatal(err)
		}
		if err := conf.Save(); err != nil {
			glog.Fatal(err)
		}
	}
	// set global configuration struct; will be accessible through config.Get()
	conf.SetCurrent()
	glog.Info("Configuration initialized")
}
