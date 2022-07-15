package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/docker/docker/client"
	"github.com/gofiber/fiber/v2"
	"github.com/yeseh/zenith/mgmt/function"
	"github.com/yeseh/zenith/mgmt/runtime"
)

func main() {
	app := fiber.New()

	db, err := sql.Open("sqlite3", "zenith-mgmt.db")
	if err != nil {
		log.Fatal(err)
	}

	repo := function.NewFunctionRepository(db)
	if err := repo.Migrate(); err != nil {
		log.Fatal(err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Zenith Management API is running!")
	})

	docker := getDockerContext()

	function.MountFunctionRouter(repo, app, docker)
	runtime.MountContainerRouter(app, docker)

	app.Listen(":3000")
}

func getDockerContext() runtime.DockerContext {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}

	return runtime.DockerContext{
		Client: cli,
		Ctx:    context.Background(),
	}
}
