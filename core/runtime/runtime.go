package runtime

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

// This file manages container instances for an application

var appContainerMap map[string]string = make(map[string]string)

func StartApp(docker DockerContext, conf StartAppConfig) (string, error) {
	containerName := fmt.Sprintf("%s-%s", conf.AppName, conf.Runtime)
	curContainer, err := FindContainer(docker, containerName)
	if err != nil {
		return "Failed to lookup containers", err
	}

	imageName := ImageNameForApp(conf.AppName, conf.Runtime)
	curContainerId := curContainer.ID
	exists := len(curContainer.ID) > 0
	isRunning := exists && curContainer.State == "running"

	// TODO: For now we support one container per app
	// If the container already exists, start if not running
	// Otherwise create, and then start
	if isRunning {
		return "Container already running", nil
	}

	if !exists {
		curContainerId, err = createContainer(docker, imageName, containerName, conf.Port)
		if err != nil {
			return "Failed to start container", err
		}
	}

	if err := docker.Client.ContainerStart(docker.Ctx, curContainerId, types.ContainerStartOptions{}); err != nil {
		return "Failed to start container", err
	}

	msg := fmt.Sprintf("Container %s started at port %s", containerName, conf.Port)

	return msg, nil
}

func createContainer(docker DockerContext, imageName string, containerName string, port string) (string, error) {
	conf := container.Config{
		Image: imageName + ":latest",
	}

	hostConfig := container.HostConfig{
		PortBindings: nat.PortMap{
			"80/tcp": []nat.PortBinding{
				{
					HostPort: port,
				},
			},
		},
	}

	resp, err := docker.Client.ContainerCreate(context.Background(), &conf, &hostConfig, nil, nil, containerName)

	return resp.ID, err
}
