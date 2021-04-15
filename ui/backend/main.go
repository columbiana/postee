package main

import (
	"log"
	"os"

	"github.com/aquasecurity/postee/ui/backend/uiserver"
)

const (
	ENV_FILELOG    = "POSTEE_UI_LOGFILE"
	ENV_CFG        = "POSTEE_UI_CFG"
	ENV_WEB        = "POSTEE_UI_WEB"
	ENV_UPDATE_URL = "POSTEE_UI_UPDATE_URL"
	ENV_UPDATE_KEY = "POSTEE_UI_KEY"
	ENV_PORT       = "POSTEE_UI_PORT"

	DEFAULT_WEB_PATH = "/uiserver/www"
)

func main() {
	logfile := os.Getenv(ENV_FILELOG)
	if logfile != "" {
		f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0444)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	cfg := os.Getenv(ENV_CFG)
	if cfg == "" {
		log.Fatalf("cfg file name is empty. You have to set a filename via %q environment variable.", ENV_CFG)
	}
	web := os.Getenv(ENV_WEB)
	if web == "" {
		web = DEFAULT_WEB_PATH
		log.Printf("The default path to web (%q) is using now.", web)
	}
	updateUrl := os.Getenv(ENV_UPDATE_URL)
	if updateUrl == "" {
		log.Printf("WARNING! Using an empty update url, UI won't restart your Postee instance with a saved configuration. You can change it via %q environment variable.", ENV_UPDATE_URL)
	}
	updateKey := os.Getenv(ENV_UPDATE_KEY)
	if updateKey == "" {
		log.Printf("WARNING! Using an empty update key. You can change it via %q environment variable.", ENV_UPDATE_KEY)
	}

	port := os.Getenv(ENV_PORT)
	if port == "" {
		port = "8090"
		log.Printf("WARNING! Using a default port: %s. You can change it via %q environment variable.", port, ENV_PORT)
	}

	server := uiserver.Instance(web, port, cfg, updateUrl, updateKey)
	server.Start()
	defer server.Stop()
}
