package routes

import (
	"github.com/ggmolly/galactf/orm"
	"github.com/ggmolly/galactf/utils"
	"github.com/gofiber/fiber/v2"
)

func GetUser(c *fiber.Ctx) error {
	if user, ok := c.Locals("user").(*orm.User); ok {
		return utils.RestStatusFactoryData(c, fiber.StatusOK, user, "")
	} else {
		return utils.RestStatusFactory(c, fiber.StatusUnauthorized, "failed to authenticate user")
	}
}
