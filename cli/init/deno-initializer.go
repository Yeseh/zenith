package init

import (
	"os"
	"path"

	zt "github.com/yeseh/zenith/template"
)

type DenoInitializer struct{}

func (d *DenoInitializer) Initialize(data InitData, output string) error {
	configOut := path.Join(output, "zenith.toml")
	if err := zt.RenderToFile(zt.ZenithTomlConfigTemplate, data, configOut); err != nil {
		return err
	}

	if err := d.createPlaceholderFunc(output); err != nil {
		return err
	}

	return nil
}

func (d *DenoInitializer) createPlaceholderFunc(basePath string) error {
	funcPath := path.Join(basePath, "functions")
	if err := os.Mkdir(funcPath, os.ModeDir); err != nil {
		return err
	}

	file, err := os.Create(path.Join(funcPath, "ping.ts"))
	if err != nil {
		return err
	}

	if err := os.WriteFile(file.Name(), []byte(zt.DenoPingFuncTemplate), 0666); err != nil {
		return err
	}

	return nil
}
