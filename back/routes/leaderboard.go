package routes

import (
	"github.com/ggmolly/galactf/orm"
	"github.com/ggmolly/galactf/utils"
	"github.com/gofiber/fiber/v2"
)

func GetLeaderboard(c *fiber.Ctx) error {
    solvedAttempts, err := orm.GetAllSolvedAttempts()

	if err != nil {
		return utils.RestStatusFactory(c, fiber.StatusInternalServerError, "error fetching solved attempts")
	}
	return utils.RestStatusFactoryData(c, fiber.StatusOK, solvedAttempts, "")
}
