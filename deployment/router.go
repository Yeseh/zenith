package deployment

import (
	"github.com/gofiber/fiber/v2"
)

func createGetAllHandler(repo *DeploymentRepository) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		results, err := repo.GetAll()
		if err != nil {
			// NOTE: this is very unlikely to be an error that the user needs to act upon
			return err
		}

		c.Status(200).JSON(results)

		return nil
	}
}

func createDeploymentHandler(repo *DeploymentRepository) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := new(CreateDeploymentDto)
		err := c.BodyParser(dto)
		if err != nil {
			c.Status(400).SendString("Invalid request!")
		}

		result, err := repo.Create(dto)
		if err != nil {
			c.Status(400).SendString("Invalid request!")
		}

		c.Status(201).JSON(result)

		return nil
	}
}

func MountDeploymentRouter(repo *DeploymentRepository, app *fiber.App) {
	root := app.Group("/deployments")

	v1 := root.Group("/v1")
	v1.Post("/", createDeploymentHandler(repo))
	v1.Get("/", createGetAllHandler(repo))
}
