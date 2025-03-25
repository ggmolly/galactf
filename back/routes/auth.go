package routes

import (
	"github.com/ggmolly/galactf/orm"
	"github.com/ggmolly/galactf/utils"
	"github.com/gofiber/fiber/v2"
)

func GetUser(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*orm.User)
	if !ok {
		return utils.RestStatusFactory(c, fiber.StatusUnauthorized, "failed to authenticate user")
	}
	return utils.RestStatusFactoryData(c, fiber.StatusOK, user, "")
}
