package middlewares

import (
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

func ReadUser(c *fiber.Ctx) *orm.User {
	return c.Locals("user").(*orm.User)
}
