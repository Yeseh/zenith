package runtime

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gofiber/fiber/v2"
)

type DockerContext struct {
	Client *client.Client
	Ctx    context.Context
}

type CreateImageResponse struct {
	Success     bool     `json:"success"`
	ImageTags   []string `json:"imageTags"`
	BuildStream string   `json:"buildStream"`
}

func listContainersHandler(docker DockerContext) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		containers, err := docker.Client.ContainerList(docker.Ctx, types.ContainerListOptions{})
		if err != nil {
			return c.SendStatus(400)
		}

		return c.Status(200).JSON(containers)
	}
}

func MountContainerRouter(app *fiber.App, docker DockerContext) {
	root := app.Group("/runtime")
	cnt := root.Group("/containers")
	// images := root.Group("/images")

	// imgv1 := images.Group("/v1")

	cntv1 := cnt.Group("/v1")
	cntv1.Get("/", listContainersHandler(docker))

	// cntv1 := containers.Group("/v1")
}
