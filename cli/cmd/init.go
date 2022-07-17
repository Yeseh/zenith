package cmd

import (
	"os"
	"path"
	"text/template"

	"github.com/spf13/cobra"
)

type InitData struct {
	Name    string
	Runtime string
}

var (
	runtime string
	output  string
	appName string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new Zenith project",
	Run: func(cmd *cobra.Command, args []string) {
		initZenith()
	},
}

var confTpl string = `[app]
name = "{{.Name}}"
runtime = "{{.Runtime}}"

[[function]]
name = "ping"
path = "./functions/ping.ts"
route = "/api/functions"
`

var pingFunc string = `export default (_: Request): Response => {
    return new Response( "Zenith is running!", {status: 200} )
};
`

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&appName, "name", "n", "my-zenith-app", "The name of your app")
	initCmd.Flags().StringVarP(&runtime, "runtime", "r", "deno", "Select the runtime for your app")
	initCmd.Flags().StringVarP(&output, "output", "o", ".", "Output path for the app")
}

func createZenithConfig(data InitData, basePath string) error {
	file, err := os.Create(path.Join(basePath, "zenith.toml"))
	if err != nil {
		return err
	}

	tmpl, err := template.New("config").Parse(confTpl)
	if err != nil {
		return err
	}

	err = tmpl.Execute(file, data)
	return nil
}

func createPlaceholderFunc(basePath string) error {
	funcPath := path.Join(basePath, "functions")
	if err := os.Mkdir(funcPath, os.ModeDir); err != nil {
		return err
	}

	file, err := os.Create(path.Join(funcPath, "ping.ts"))
	if err != nil {
		return err
	}

	if err := os.WriteFile(file.Name(), []byte(pingFunc), 0666); err != nil {
		return err
	}

	return nil
}

func initZenith() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	config := InitData{
		Name:    appName,
		Runtime: runtime,
	}

	if err := createZenithConfig(config, wd); err != nil {
		return err
	}

	if err := createPlaceholderFunc(wd); err != nil {
		return err
	}

	return nil
}
