package runtime

import (
	"github.com/docker/docker/api/types"
	"github.com/gofiber/fiber/v2"
)

func listContainersHandler(docker DockerContext) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		containers, err := docker.Client.ContainerList(docker.Ctx, types.ContainerListOptions{})
		if err != nil {
			return c.SendStatus(400)
		}

		return c.Status(200).JSON(containers)
	}
}

// Lists all zenith-app images currently in the registry
func listImagesHandler(docker DockerContext) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		images, err := ListImages(docker)
		if err != nil {
			return c.SendStatus(400)
		}

		return c.Status(200).JSON(images)
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
