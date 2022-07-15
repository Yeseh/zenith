package runtime

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

var appContainerMap map[string]string = make(map[string]string)

type StartAppResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	ContainerID string `json:"containerId"`
	Image       string `json:"image"`
}

func StartApp(appName string, runtime string, docker DockerContext) (string, error) {
	port := "5000"
	containers, err := docker.Client.ContainerList(docker.Ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return "", err
	}

	imageName := ImageNameForApp(appName, runtime)
	containerName := fmt.Sprintf("%s-%s", appName, runtime)

	var curContainer types.Container = types.Container{}
	var curContainerId string

	for _, cnt := range containers {
		exists := cnt.Names[0] == "/"+containerName
		if exists {
			curContainer = cnt
			break
		}
	}

	if len(curContainer.ID) > 0 {
		curContainerId = curContainer.ID
		if curContainer.State == "running" {
			return "Container already running", nil
		}
	} else {
		curContainerId, err = createContainer(docker, imageName, containerName, port)
		if err != nil {
			return "Failed to start container", err
		}
	}

	if err := docker.Client.ContainerStart(docker.Ctx, curContainerId, types.ContainerStartOptions{}); err != nil {
		return "Failed to start container", err
	}

	msg := fmt.Sprintf("Container %s started at port %s", containerName, port)
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
