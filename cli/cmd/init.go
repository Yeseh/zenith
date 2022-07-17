package cmd

import (
	"os"

	"github.com/spf13/cobra"
	zci "github.com/yeseh/zenith/cli/init"
)

var (
	runtime    string
	outputPath string
	appName    string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new Zenith project",
	Run: func(cmd *cobra.Command, args []string) {
		wd := outputPath
		if len(wd) == 0 {
			oswd, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			wd = oswd
		}

		i, err := zci.InitializerFactory(runtime)
		if err != nil {
			panic(err)
		}

		data := zci.InitData{
			AppName: appName,
			Runtime: runtime,
		}

		if err := i.Initialize(data, outputPath); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&appName, "name", "n", "my-zenith-app", "The name of your app")
	initCmd.Flags().StringVarP(&runtime, "runtime", "r", "deno", "Select the runtime for your app")
	initCmd.Flags().StringVarP(&outputPath, "output", "o", ".", "Output path for the app")
}
