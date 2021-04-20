package main

import (
	"fmt"
	"github.com/aquasecurity/postee/alertmgr"
	"github.com/aquasecurity/postee/utils"
	"github.com/aquasecurity/postee/webserver"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
)

const (
	URL        = "0.0.0.0:8082"
	TLS        = "0.0.0.0:8445"
	URL_USAGE  = "The socket to bind to, specified using host:port."
	TLS_USAGE  = "The TLS socket to bind to, specified using host:port."
	CFG_USAGE  = "The folder which contains alert configuration files."
	CFG_FOLDER = "/config/"
)

var (
	url       = ""
	tls       = ""
	cfgFolder = ""
)

var rootCmd = &cobra.Command{
	Use:   "webhooksrv",
	Short: fmt.Sprintf("Aqua Container Security Webhook server\n"),
	Long:  fmt.Sprintf("Aqua Container Security Webhook server\n"),
}

func init() {
	rootCmd.Flags().StringVar(&url, "url", URL, URL_USAGE)
	rootCmd.Flags().StringVar(&tls, "tls", TLS, TLS_USAGE)
	rootCmd.Flags().StringVar(&cfgFolder, "cfgFolder", CFG_FOLDER, CFG_USAGE)
}

func getFilesFromFolder(folder string) ([]string, error) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	names := []string{}
	if !strings.HasSuffix(folder, "/") {
		folder += "/"
	}
	for _, file := range files {
		names = append(names, folder+file.Name())
	}
	return names, nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	utils.InitDebug()

	rootCmd.Run = func(cmd *cobra.Command, args []string) {

		if os.Getenv("AQUAALERT_URL") != "" {
			url = os.Getenv("AQUAALERT_URL")
		}

		if os.Getenv("AQUAALERT_TLS") != "" {
			tls = os.Getenv("AQUAALERT_TLS")
		}

		cfgFolder = os.Getenv("AQUAALERT_CFG_FOLDER")
		if cfgFolder == "" {
			log.Printf("AQUAALERT_CFG_FOLDER environment variable is empty!")
			return
		}
		files, err := getFilesFromFolder(cfgFolder)
		if err != nil {
			log.Printf("getFilesFromFolder(%q) error: %v", cfgFolder, err)
			return
		}

		go alertmgr.Instance().Start(files)
		defer alertmgr.Instance().Terminate()

		go webserver.Instance().Start(url, tls)
		defer webserver.Instance().Terminate()

		Daemonize()
	}
	rootCmd.Execute()
}

func Daemonize() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Println(sig)
		done <- true
	}()

	<-done
}
