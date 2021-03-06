package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"runtime"

	"github.com/drud/ddev/pkg/appports"
	"github.com/drud/ddev/pkg/plugins/platform"
	"github.com/drud/ddev/pkg/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// SequelproLoc is where we expect to find the sequel pro.app
// It's global so it can be mocked in testing.
var SequelproLoc = "/Applications/sequel pro.app"

// localDevSequelproCmd represents the sequelpro command
var localDevSequelproCmd = &cobra.Command{
	Use:   "sequelpro",
	Short: "Easily connect local site to sequelpro",
	Long:  `A helper command for easily using sequelpro (OSX database browser) with a ddev app that has been initialized locally.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			log.Fatalf("invalid arguments to sequelpro command: %v", args)
		}

		out, err := handleSequelProCommand(SequelproLoc)
		if err != nil {
			log.Fatalf("Could not run sequelpro command: %s", err)
		}
		util.Success(out)
	},
}

// handleSequelProCommand() is the "real" handler for the real command
func handleSequelProCommand(appLocation string) (string, error) {
	app, err := getActiveApp()
	if err != nil {
		return "", err
	}

	if app.SiteStatus() != "running" {
		return "", errors.New("app not running locally. Try `ddev start`")
	}

	dbPort := appports.GetPort("db")

	tmpFilePath := filepath.Join(app.AppRoot(), ".ddev/sequelpro.spf")
	tmpFile, err := os.Create(tmpFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer util.CheckClose(tmpFile)

	_, err = tmpFile.WriteString(fmt.Sprintf(
		platform.SequelproTemplate,
		"data",         //dbname
		app.HostName(), //host
		app.HostName(), //connection name
		"root",         // dbpass
		dbPort,         // port
		"root",         //dbuser
	))
	util.CheckErr(err)

	err = exec.Command("open", tmpFilePath).Run()
	if err != nil {
		return "", err
	}
	return "sequelpro command finished successfully!", nil
}

// dummyDevSequelproCmd represents the "not available" sequelpro command
var dummyDevSequelproCmd = &cobra.Command{
	Use:   "sequelpro",
	Short: "This command is not available since sequel pro.app is not installed",
	Long:  `Where installed, "ddev sequelpro" launches the sequel pro database browser`,
	Run: func(cmd *cobra.Command, args []string) {
		util.Failed("The sequelpro command is not available because sequel pro.app is not detected on your workstation")

	},
}

// init installs the real command if it's available, otherwise dummy command (if on OSX), otherwise no command
func init() {
	switch {
	case detectSequelpro():
		RootCmd.AddCommand(localDevSequelproCmd)
	case runtime.GOOS == "darwin":
		RootCmd.AddCommand(dummyDevSequelproCmd)
	}
}

// detectSequelpro looks for the sequel pro app in /Applications; returns true if found
func detectSequelpro() bool {
	if _, err := os.Stat(SequelproLoc); err == nil {
		return true
	}
	return false
}
