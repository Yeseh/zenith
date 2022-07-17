package runtime

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type DockerContext struct {
	Client *client.Client
	Ctx    context.Context
}

type ImageConfig struct {
	Runtime            string
	AppName            string
	Location           string
	DockerBuildOptions types.ImageBuildOptions
}

type StartAppConfig struct {
	AppName       string
	Runtime       string
	Port          string
	InstanceCount int
}

// TODO: JSON types -> API
type StartAppResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	ContainerID string `json:"containerId"`
	Image       string `json:"image"`
}

type CreateImageResponse struct {
	Success     bool     `json:"success"`
	ImageTags   []string `json:"imageTags"`
	BuildStream string   `json:"buildStream"`
}
