package runtime

import (
	"github.com/docker/docker/api/types"
)

func FindContainer(docker DockerContext, containerName string) (types.Container, error) {
	opts := types.ContainerListOptions{All: true}
	containers, err := docker.Client.ContainerList(docker.Ctx, opts)
	if err != nil {
		return types.Container{}, err
	}

	var curContainer types.Container = types.Container{}

	for _, cnt := range containers {
		exists := cnt.Names[0] == "/"+containerName
		if exists {
			curContainer = cnt
			break
		}
	}

	return curContainer, nil
}
