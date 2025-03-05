package routes

import (
	"log"

	"github.com/ggmolly/galactf/orm"
	"github.com/ggmolly/galactf/utils"
	"github.com/gofiber/fiber/v2"
)

func GetChallenges(c *fiber.Ctx) error {
	chals, err := orm.GetChallengeStats()
	status := fiber.StatusOK
	if err != nil {
		log.Printf("[-] error fetching challenges: %s", err.Error())
		status = fiber.StatusInternalServerError
	}
	return utils.RestStatusFactoryData(c, status, chals, "")
}

func GetChallenge(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.RestStatusFactory(c, fiber.StatusBadRequest, "invalid id")
	}
	chal, err := orm.GetChallengeStatsById(id)
	if err != nil {
		return utils.RestStatusFactory(c, fiber.StatusInternalServerError, "error fetching challenge")
	}
	return utils.RestStatusFactoryData(c, fiber.StatusOK, chal, "")
}

func GetSolvers(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.RestStatusFactory(c, fiber.StatusBadRequest, "invalid id")
	}
	attempts, err := orm.GetSolvedAttempts(id)
	if err != nil {
		return utils.RestStatusFactory(c, fiber.StatusInternalServerError, "error fetching challenge")
	}
	return utils.RestStatusFactoryData(c, fiber.StatusOK, attempts, "")
}
