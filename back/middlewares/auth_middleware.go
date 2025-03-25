package middlewares

import (
	"log"
	"os"

	"github.com/ggmolly/galactf/orm"
	"github.com/ggmolly/galactf/utils"
	"github.com/gofiber/fiber/v2"
)

var (
	agnosticMiddleware func(c *fiber.Ctx) error = nil
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

func EnvAuthMiddleware(c *fiber.Ctx) func(c *fiber.Ctx) error {
	if agnosticMiddleware == nil && os.Getenv("MODE") == "dev" {
		agnosticMiddleware = DummyAuthMiddleware
	} else if agnosticMiddleware == nil && os.Getenv("MODE") == "prod" {
		agnosticMiddleware = GaladrimAuthMiddleware
	} else {
		log.Fatal("invalid mode:", os.Getenv("MODE"))
	}
	return agnosticMiddleware
}

func ReadUser(c *fiber.Ctx) *orm.User {
	return c.Locals("user").(*orm.User)
}
