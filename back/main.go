package main

import (
	"log"
	"os"

	"github.com/bytedance/sonic"
	"github.com/ggmolly/galactf/orm"
	"github.com/ggmolly/galactf/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file, does it exist?")
	}

	orm.InitDatabase()

	if len(os.Args) > 1 && os.Args[1] == "seed" {
		orm.Seed()
	}
}

func main() {
	app := fiber.New(fiber.Config{
		AppName:     "galactf",
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	apiGroup := app.Group("/api/v1")
	{
		apiGroup.Get("/challenges", routes.GetChallenges)
	}

	app.Listen(":7777")
}
