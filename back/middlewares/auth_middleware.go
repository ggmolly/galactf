package middlewares

import (
	"log"
	"os"

	"github.com/ggmolly/galactf/orm"
	"github.com/ggmolly/galactf/utils"
	"github.com/gofiber/fiber/v2"
)

func DummyAuthMiddleware(c *fiber.Ctx) error {
	var user orm.User
	if err := orm.GormDB.First(&user).Error; err != nil {
		return utils.RestStatusFactory(c, fiber.StatusUnauthorized, "failed to authenticate user")
	}
	c.Locals("user", &user)
	return c.Next()
}

func GaladrimAuthMiddleware(c *fiber.Ctx) error {
	user, err := orm.GetUserFromCookie(c)
	if err != nil {
		return utils.RestStatusFactory(c, fiber.StatusUnauthorized, "failed to authenticate user")
	}
	c.Locals("user", user)
	return c.Next()
}

func AgnosticAuthMiddleware() func(c *fiber.Ctx) error {
	if os.Getenv("MODE") == "dev" {
		log.Println("[!] using dummy auth middleware in dev mode")
		return DummyAuthMiddleware
	} else if os.Getenv("MODE") == "prod" {
		log.Println("[!] using galadrim auth middleware in prod mode")
		return GaladrimAuthMiddleware
	} else {
		log.Fatal("invalid mode:", os.Getenv("MODE"))
	}
	// dead code
	return DummyAuthMiddleware
}

func ReadUser(c *fiber.Ctx) *orm.User {
	return c.Locals("user").(*orm.User)
}
