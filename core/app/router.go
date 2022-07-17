package app

// import (
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/yeseh/zenith/runtime"
// 	zr "github.com/yeseh/zenith/runtime"
// )

// func createGetAllHandler(repo *AppRepository) func(*fiber.Ctx) error {
// 	return func(c *fiber.Ctx) error {
// 		results, err := repo.GetAll()
// 		if err != nil {
// 			return err
// 		}

// 		if len(results) == 0 {
// 			return c.SendStatus(204)
// 		}

// 		return c.Status(200).JSON(results)
// 	}
// }

// func createGetByNameHandler(repo *AppRepository) func(*fiber.Ctx) error {
// 	return func(c *fiber.Ctx) error {
// 		name := c.Params("displayName")
// 		if len(name) == 0 {
// 			return c.SendStatus(404)
// 		}

// 		results, err := repo.GetByDisplayName(name)
// 		if err != nil {
// 			return err
// 		}

// 		return c.Status(200).JSON(results)
// 	}
// }

// func createStartAppHandler(repo *AppRepository, docker zr.DockerContext) func(*fiber.Ctx) error {
// 	return func(c *fiber.Ctx) error {
// 		name := c.Params("displayName")

// 		if len(name) == 0 {
// 			return c.SendStatus(404)
// 		}

// 		results, err := repo.GetByDisplayName(name)
// 		if err != nil {
// 			return c.SendStatus(400)
// 		}

// 		if len(results.AppName) == 0 {
// 			return c.SendStatus(404)
// 		}

// 		conf := runtime.StartAppConfig{
// 			AppName: results.AppName,
// 			Runtime: results.Runtime,
// 			// TODO: Port should probably be managed by the runtime
// 			Port: "5000",
// 		}

// 		msg, err := zr.StartApp(docker, conf)
// 		if err != nil {
// 			return c.Status(400).SendString(err.Error())
// 		}

// 		return c.Status(202).SendString(msg)
// 	}
// }

// func createAppHandler(repo *AppRepository, docker zr.DockerContext) func(*fiber.Ctx) error {
// 	return func(c *fiber.Ctx) error {
// 		dto := new(CreateAppDto)

// 		if err := c.BodyParser(dto); err != nil {
// 			return c.Status(400).SendString("Invalid request!")
// 		}

// 		// TODO: validate runtimes/definition names
// 		if dto.Runtime != "deno" {
// 			return c.Status(404).SendString("Unsupported runtime")
// 		}

// 		function, err := CreateApp(docker, repo, dto)
// 		if err != nil {
// 			return c.SendStatus(400)
// 		}

// 		return c.Status(201).JSON(function)
// 	}
// }

// func MountAppRouter(repo *AppRepository, app *fiber.App, docker zr.DockerContext) {
// 	root := app.Group("/app")

// 	v1 := root.Group("/v1")
// 	v1.Post("/", createAppHandler(repo, docker))
// 	v1.Get("/:displayName", createGetByNameHandler(repo))
// 	v1.Put("/:displayName/start", createStartAppHandler(repo, docker))
// 	// v1.Get("/:displayName/stop", createGetByNameHandler(repo))
// 	v1.Get("/", createGetAllHandler(repo))
// }
