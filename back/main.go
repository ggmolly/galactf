package main

import (
	"log"
	"os"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/ggmolly/galactf/factories"
	"github.com/ggmolly/galactf/middlewares"
	"github.com/ggmolly/galactf/orm"
	"github.com/ggmolly/galactf/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
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

	app.Use(cors.New(cors.Config{
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowOriginsFunc: func(origin string) bool {
			return strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1")
		},
		AllowCredentials: true,
	}))

	app.Use(logger.New(logger.Config{
		Format: "[${time}] [${ip}] [${method}] [${status}] @ ${path} | ${latency}\n",
		TimeZone: "UTC",
		DisableColors: true,
	}))

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			// TODO: discord webhook
			log.Println(e)
		},
	}))

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
				oneTrickGroup.Post("/", middlewares.ChallengeUnlockedMiddleware("one trick"), factories.SubmitOneTrick)
				oneTrickGroup.Post("/encrypt", middlewares.ChallengeUnlockedMiddleware("one trick"), factories.EncryptOneTrick)
			}

		}
	}

	app.Listen(":7777")
}
