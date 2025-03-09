package routes

import (
	"github.com/ggmolly/galactf/orm"
	"github.com/ggmolly/galactf/utils"
	"github.com/gofiber/fiber/v2"
)

func GetUser(c *fiber.Ctx) error {
	// TODO: Implement auth middleware
	var user orm.User
	if err := orm.GormDB.First(&user).Error; err != nil {
		return utils.RestStatusFactory(c, fiber.StatusUnauthorized, "failed to authenticate user")
	}
	return utils.RestStatusFactoryData(c, fiber.StatusOK, user, "")
}
