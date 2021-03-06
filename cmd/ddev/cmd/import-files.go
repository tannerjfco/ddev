package cmd

import (
	"log"
	"os"

	"github.com/drud/ddev/pkg/util"
	"github.com/spf13/cobra"
)

var fileSource string

// ImportFileCmd represents the `ddev import-db` command.
var ImportFileCmd = &cobra.Command{
	Use:   "import-files",
	Short: "Import the uploaded files directory of an existing site to the default public upload directory of your application.",
	Long:  "Import the uploaded files directory of an existing site to the default public upload directory of your application. The files can be provided as a directory path or an archive in .tar.gz format. The contents at the root of the archive or directory will be the contents of the default public upload directory of your application. If the destination directory exists, it will be replaced with the assets being imported.",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			err := cmd.Usage()
			util.CheckErr(err)
			os.Exit(0)
		}

		client := util.GetDockerClient()

		err := util.EnsureNetwork(client, netName)
		if err != nil {
			log.Fatal(err)
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		app, err := getActiveApp()
		if err != nil {
			util.Failed("Failed to import files: %v", err)
		}

		err = app.ImportFiles(fileSource)
		if err != nil {
			util.Failed("Failed to import files for %s: %v", app.GetName(), err)
		}
		util.Success("Successfully imported files for %v", app.GetName())
	},
}

func init() {
	ImportFileCmd.Flags().StringVarP(&fileSource, "src", "", "", "Provide the path to a directory or .tar.gz archive of files to import")
	RootCmd.AddCommand(ImportFileCmd)
}
