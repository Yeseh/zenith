package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/yeseh/zenith/mgmt/deployment"
)

func main() {
	app := fiber.New()

	db, err := sql.Open("sqlite3", "zenith.db")
	if err != nil {
		log.Fatal(err)
	}

	repo := deployment.NewDeploymentRepository(db)
	if err := repo.Migrate(); err != nil {
		log.Fatal(err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API is running!")
	})

	deployment.MountDeploymentRouter(repo, app)

	app.Listen(":3000")
}
