package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/ggmolly/galactf/cache"
	"github.com/ggmolly/galactf/factories"
	"github.com/ggmolly/galactf/middlewares"
	"github.com/ggmolly/galactf/orm"
	"github.com/ggmolly/galactf/routes"
	"github.com/ggmolly/galactf/types"
	"github.com/ggmolly/galactf/utils"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/google/uuid"
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
	cache.InitRedisClient()
	utils.InitRedisStore()
	go routes.RevealAgent()

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
		ProxyHeader: "X-Forwarded-For",
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("name", "anonymous")
		return c.Next()
	})

	app.Use(cors.New(cors.Config{
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowOriginsFunc: func(origin string) bool {
			return strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1")
		},
		AllowCredentials: true,
	}))

	app.Use(logger.New(logger.Config{
		Format:        "[${time}] [${ip} / ${locals:name}] [${method}] [${status}] @ ${path} | ${latency}\n",
		TimeZone:      "UTC",
		DisableColors: true,
	}))

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			buf := make([]byte, 1<<16)
			n := runtime.Stack(buf, true)
			
			fileName := uuid.NewString() + ".dmp"
			fileFail := false

			f, err := os.Create(fileName)
			if err != nil {
				log.Printf("[!] Failed to create stacktrace dump file: %s", err.Error())
				fileFail = true
			} else {
				defer f.Close()
			}

			buf = buf[:n]

			_, err = f.Write(buf)
			if err != nil {
				log.Printf("[!] Failed to write stacktrace dump file: %s", err.Error())
				fileFail = true
			}

			log.Printf("[!] Panic: %s\n%s", e, buf)

			message := types.Message{
				Blocks: []types.Block{
					{
						Type: "header",
						Text: &types.Text{
							Type: "plain_text",
							Text: ":rotating_light: Panic",
						},
					},
					{
						Type: "section",
						Text: &types.Text{
							Type: "mrkdwn",
							Text: "Détails: `" + fmt.Sprint(e) + "`",
						},
					},
					{
						Type: "section",
						Text: &types.Text{
							Type: "mrkdwn",
							Text: "Stacktrace dump: `" + fileName + "`",
						},
					},
					{
						Type: "section",
						Text: &types.Text{
							Type: "mrkdwn",
							Text: fmt.Sprintf("Flag fileFail: `%v`", fileFail),
						},
					},
				},
			}

			

			types.SendSlackWebhook(os.Getenv("PANIC_SLACK_WEBHOOK_URI"), &message)
		},
	}))

	authMiddleware := middlewares.AgnosticAuthMiddleware()

	apiGroup := app.Group("/api/v1")
	{
		apiGroup.Get("/ws", authMiddleware, middlewares.WsUpgradeMiddleware, websocket.New(routes.WsHandler))

		authGroup := apiGroup.Group("/auth", authMiddleware)
		{
			authGroup.Get("/me", routes.GetUser)
		}
		challengesGroup := apiGroup.Group("/challenges", authMiddleware)
		{
			challengesGroup.Get("/", routes.GetChallenges)
		}
        leaderboardGroup := apiGroup.Group("/leaderboard", authMiddleware)
        {
            leaderboardGroup.Get("/", routes.GetLeaderboard)
        }
		challengeGroup := apiGroup.Group("/challenge/:id", authMiddleware)
		{
			challengeGroup.Get("/", routes.GetChallenge)
			challengeGroup.Get("/solvers", routes.GetSolvers)
			challengeGroup.Post("/submit", routes.SubmitFlag)
		}
		factoriesGroup := apiGroup.Group("/factories", authMiddleware)
		{
			factoriesGroup.Get("/elite_encryption", middlewares.ChallengeUnlockedMiddleware("elite encryption"), factories.GenerateEliteEncryption)
			factoriesGroup.Get("/super_elite_encryption", middlewares.ChallengeUnlockedMiddleware("super elite encryption"), factories.GenerateSuperEliteEncryption)
			factoriesGroup.Get("/quack", middlewares.ChallengeUnlockedMiddleware("quack"), factories.GenerateQuackChallenge)
			factoriesGroup.Get("/more_or_less", middlewares.ChallengeUnlockedMiddleware("more or less"), factories.GenerateMoreOrLess)
			factoriesGroup.Get("/cat_image", middlewares.ChallengeUnlockedMiddleware("cat image"), factories.GenerateCatImage)
			factoriesGroup.Get("/quiet_riot_code", middlewares.ChallengeUnlockedMiddleware("quiet riot code"), factories.GenerateQuietRiotCode)
			factoriesGroup.Get("/byte_bounty",
				middlewares.ChallengeUnlockedMiddleware("byte bounty"),
				middlewares.NewRateLimiterMiddleware(1, 30*time.Minute, "HTTP 429: 30 minutes are required between two calls to this endpoint."),
				factories.GenerateByteBounty,
			)

			// "Exclusive club" challenge
			oneTrickGroup := factoriesGroup.Group("/exclusive_club")
			{
				oneTrickGroup.Get("/", middlewares.ChallengeUnlockedMiddleware("exclusive club"), factories.RenderOneTrick)
				oneTrickGroup.Post("/", middlewares.ChallengeUnlockedMiddleware("exclusive club"), factories.SubmitOneTrick)
				oneTrickGroup.Post("/encrypt", middlewares.ChallengeUnlockedMiddleware("exclusive club"), factories.EncryptOneTrick)
			}

			// Every proxied challenges, the User ID is passed in 'X-User-ID' header
			// and the flag matching the user is passed in 'X-GalaCTF-Flag' header
			// 1st argument is: the container name in docker-compose.yml
			// 2nd argument is: the base URL of the challenge (which is the URL that's in the Group() call)
			// 3rd argument is: the name of the challenge in the database (used to generate & validate the flag)
			factoriesGroup.Group("/bobby_library", middlewares.ChallengeUnlockedMiddleware("bobby's library"), routes.ProxyFactory("bobby_library", "/bobby_library", "bobby's library"))
			factoriesGroup.Group("/unsecure_notes", middlewares.ChallengeUnlockedMiddleware("(un)secure notes"), routes.ProxyFactory("unsecure_notes", "/unsecure_notes", "(un)secure notes"))
			factoriesGroup.Group("/just_a_moment", middlewares.ChallengeUnlockedMiddleware("just a moment..."), routes.ProxyFactory("just_a_moment", "/just_a_moment", "just a moment..."))
			factoriesGroup.Group("/cookie_monster_squared", middlewares.ChallengeUnlockedMiddleware("cookie monster²"), routes.ProxyFactory("cookie_monster_squared", "/cookie_monster_squared", "cookie monster²"))
			factoriesGroup.Group("/cookie_monster", middlewares.ChallengeUnlockedMiddleware("cookie monster"), routes.ProxyFactory("cookie_monster", "/cookie_monster", "cookie monster"))
			factoriesGroup.Group("/calculator", middlewares.ChallengeUnlockedMiddleware("calculator"), routes.ProxyFactory("calculator", "/calculator", "calculator"))
			factoriesGroup.Group("/claustrophobia", middlewares.ChallengeUnlockedMiddleware("claustrophobia"), routes.ProxyFactory("claustrophobia", "/claustrophobia", "claustrophobia"))
			factoriesGroup.Group("/twisted_mersenne", middlewares.ChallengeUnlockedMiddleware("twisted mersenne"), routes.ProxyFactory("twisted_mersenne", "/twisted_mersenne", "twisted mersenne"))
			factoriesGroup.Group("/regex_battle", middlewares.ChallengeUnlockedMiddleware("anti-spirit team"), routes.ProxyFactory("regex_battle", "/regex_battle", "anti-spirit team"))
		}
	}

	app.Listen("0.0.0.0:7777")
}
