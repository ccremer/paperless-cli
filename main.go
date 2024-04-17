package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ccremer/paperless-cli/pkg/app"
	"github.com/ccremer/plogr"
)

var (
	// These will be populated by Goreleaser
	version = "unknown"
	commit  = "-dirty-"
	date    = time.Now().Format("2006-01-02")
)

func main() {
	appInstance := app.NewApp()
	appInstance.Version = fmt.Sprintf("%s, revision=%s, date=%s", version, commit, date)
	err := appInstance.Run(os.Args)
	if err != nil {
		plogr.DefaultErrorPrinter.Println(err.Error())
		os.Exit(1)
	}
}
