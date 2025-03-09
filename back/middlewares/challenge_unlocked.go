package middlewares

import (
	"time"

	"github.com/ggmolly/galactf/orm"
	"github.com/ggmolly/galactf/utils"
	"github.com/gofiber/fiber/v2"
)

func ChallengeUnlockedMiddleware(name string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var chal orm.Challenge
		if err := orm.GormDB.Where("name = ?", name).First(&chal).Error; err != nil {
			return utils.RestStatusFactory(c, fiber.StatusInternalServerError, "error fetching challenge")
		}

		// Check if the challenge is locked
		if chal.RevealAt.After(time.Now().UTC()) {
			return utils.RestStatusFactory(c, fiber.StatusForbidden, "This challenge is locked")
		}

		return c.Next()
	}
}
