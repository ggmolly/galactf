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
	"github.com/gofiber/contrib/websocket"
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

	if os.Getenv("MODE") != "dev" && os.Getenv("MODE") != "prod" {
		log.Fatal("MODE must be either 'dev' or 'prod', check your .env file")
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
		Format:        "[${time}] [${ip}] [${method}] [${status}] @ ${path} | ${latency}\n",
		TimeZone:      "UTC",
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
		apiGroup.Get("/ws", middlewares.DummyAuthMiddleware, middlewares.WsUpgradeMiddleware, websocket.New(routes.WsHandler))

		authGroup := apiGroup.Group("/auth", middlewares.DummyAuthMiddleware)
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
			factoriesGroup.Get("/super_elite_encryption", middlewares.ChallengeUnlockedMiddleware("super elite encryption"), factories.GenerateSuperEliteEncryption)
			factoriesGroup.Get("/quack", middlewares.ChallengeUnlockedMiddleware("quack"), factories.GenerateQuackChallenge)
			factoriesGroup.Get("/more_or_less", middlewares.ChallengeUnlockedMiddleware("more or less"), factories.GenerateMoreOrLess)
			factoriesGroup.Get("/cat_image", middlewares.ChallengeUnlockedMiddleware("cat image"), factories.GenerateCatImage)
			factoriesGroup.Get("/quiet_riot_code", middlewares.ChallengeUnlockedMiddleware("quiet riot code"), factories.GenerateQuietRiotCode)

			// "One trick" challenge
			oneTrickGroup := factoriesGroup.Group("/one_trick")
			{
				oneTrickGroup.Get("/", middlewares.ChallengeUnlockedMiddleware("one trick"), factories.RenderOneTrick)
				oneTrickGroup.Post("/", middlewares.ChallengeUnlockedMiddleware("one trick"), factories.SubmitOneTrick)
				oneTrickGroup.Post("/encrypt", middlewares.ChallengeUnlockedMiddleware("one trick"), factories.EncryptOneTrick)
			}

			// Every proxied challenges, the User ID is passed in 'X-User-ID' header
			// and the flag matching the user is passed in 'X-GalaCTF-Flag' header
			// 1st argument is: the container name in docker-compose.yml
			// 2nd argument is: the base URL of the challenge (which is the URL that's in the Group() call)
			// 3rd argument is: the name of the challenge in the database (used to generate & validate the flag)
			factoriesGroup.Group("/bobby_library", middlewares.ChallengeUnlockedMiddleware("bobby's library"), routes.ProxyFactory("bobby_library", "/bobby_library", "bobby's library"))
			factoriesGroup.Group("/unsecure_notes", middlewares.ChallengeUnlockedMiddleware("(un)secure notes"), routes.ProxyFactory("unsecure_notes", "/unsecure_notes", "(un)secure notes"))
			factoriesGroup.Group("/just_a_moment", middlewares.ChallengeUnlockedMiddleware("just a moment..."), routes.ProxyFactory("just_a_moment", "/just_a_moment", "just a moment..."))
			factoriesGroup.Group("/cookie_monster", middlewares.ChallengeUnlockedMiddleware("cookie monster"), routes.ProxyFactory("cookie_monster", "/cookie_monster", "cookie monster"))
		}
	}

	app.Listen(":7777")
}
