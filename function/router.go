package function

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	zr "github.com/yeseh/zenith/mgmt/runtime"
	zs "github.com/yeseh/zenith/mgmt/storage"
)

func createGetAllHandler(repo *FunctionRepository) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		results, err := repo.GetAll()
		if err != nil {
			return err
		}

		if len(results) == 0 {
			return c.SendStatus(204)
		}

		return c.Status(200).JSON(results)
	}
}

func createGetByNameHandler(repo *FunctionRepository) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		name := c.Params("displayName")
		if len(name) == 0 {
			return c.SendStatus(404)
		}

		results, err := repo.GetByDisplayName(name)
		if err != nil {
			return err
		}

		return c.Status(200).JSON(results)
	}
}

type StartFunctionResponse struct {
	Success bool
	Address string
}

func createStartFunctionHandler(repo *FunctionRepository, docker zr.DockerContext) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		name := c.Params("displayName")
		if len(name) == 0 {
			return c.SendStatus(404)
		}

		results, err := repo.GetByDisplayName(name)
		if err != nil {
			return c.SendStatus(400)
		}

		if len(results.AppName) == 0 {
			return c.SendStatus(404)
		}

		msg, err := zr.StartApp(results.AppName, results.Runtime, docker)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		return c.Status(202).SendString(msg)
	}
}

// TODO: this does a lot, needs refactoring
func createFunctionHandler(repo *FunctionRepository, docker zr.DockerContext) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := new(CreateFunctionDto)

		if err := c.BodyParser(dto); err != nil {
			return c.Status(400).SendString("Invalid request!")
		}

		// TODO: validate runtimes/definition names
		if dto.Runtime != "deno" {
			return c.Status(404).SendString("Unsupported runtime")
		}

		fmt.Println("Uploading app from: ", dto.SourceLocation)
		target, err := zs.UploadApp(dto.SourceLocation, dto.AppName)
		if err != nil {
			log.Fatal(err)
		}

		config := zr.CreateImageConfig(dto.Runtime, dto.AppName, target)
		imageResponse, err := zr.CreateImageFor(config, docker)
		if err != nil || !imageResponse.Success {
			log.Fatal(err)
		}

		function := Function{
			AppName:  dto.AppName,
			Runtime:  dto.Runtime,
			Location: target,
			Images:   config.DockerBuildOptions.Tags,
		}
		result, err := repo.Create(function)

		if err != nil {
			return c.Status(400).SendString("Invalid request!")
		}

		return c.Status(201).JSON(result)
	}
}

func MountFunctionRouter(repo *FunctionRepository, app *fiber.App, docker zr.DockerContext) {
	root := app.Group("/functions")

	v1 := root.Group("/v1")
	v1.Post("/", createFunctionHandler(repo, docker))
	v1.Get("/:displayName", createGetByNameHandler(repo))
	v1.Put("/:displayName/start", createStartFunctionHandler(repo, docker))
	// v1.Get("/:displayName/stop", createGetByNameHandler(repo))
	v1.Get("/", createGetAllHandler(repo))
}
