package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var configStr = `
[app]
name = "test"
runtime = "deno"

[[functions]]
name = "ping"
path = "./functions/ping.ts"
route = "/api/functions"
`

func createTestingConfig(config string) *os.File {
	configFile, err := os.CreateTemp("", ".zenithtestconfig")
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(configFile.Name(), []byte(config), 0666); err != nil {
		panic(err)
	}

	return configFile
}

func TestReadConfig(t *testing.T) {
	file := createTestingConfig(configStr)
	defer os.Remove(file.Name())

	config, err := ReadConfig(file.Name())
	if err != nil {
		panic(err)
	}

	assert.Equal(t, "test", config.App.Name)
	assert.Equal(t, "deno", config.App.Runtime)
	assert.Len(t, config.Functions, 1)

	fn := config.Functions[0]
	assert.Equal(t, fn.Name, "ping")
	assert.Equal(t, fn.Path, "./functions/ping.ts")
	assert.Equal(t, fn.Route, "/api/functions")
}
