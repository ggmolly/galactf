package main

import (
	"log"
	"os"

	"github.com/bytedance/sonic"
	"github.com/ggmolly/galactf/factories"
	"github.com/ggmolly/galactf/middlewares"
	"github.com/ggmolly/galactf/orm"
	"github.com/ggmolly/galactf/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
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
	// templating engine (used for challenges' front-ends)
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		AppName:     "galactf",
		Views:       engine,
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	apiGroup := app.Group("/api/v1")
	{
		authGroup := apiGroup.Group("/auth")
		{
			authGroup.Get("/me", routes.GetUser)
		}
		challengesGroup := apiGroup.Group("/challenges", middlewares.DummyAuthMiddleware)
		{
			challengesGroup.Get("/", routes.GetChallenges)
		}
		challengeGroup := apiGroup.Group("/challenge/:id", middlewares.DummyAuthMiddleware)
		{
			challengeGroup.Get("/", routes.GetChallenge)
			challengeGroup.Get("/solvers", routes.GetSolvers)
			challengeGroup.Post("/submit", routes.SubmitFlag)
		}
		factoriesGroup := apiGroup.Group("/factories", middlewares.DummyAuthMiddleware)
		{
			factoriesGroup.Get("/elite_encryption", middlewares.ChallengeUnlockedMiddleware("elite encryption"), factories.GenerateEliteEncryption)

			// "One trick" challenge
			oneTrickGroup := factoriesGroup.Group("/one_trick")
			{
				oneTrickGroup.Get("/", middlewares.ChallengeUnlockedMiddleware("one trick"), factories.RenderOneTrick)
			}

		}
	}

	app.Listen(":7777")
}
