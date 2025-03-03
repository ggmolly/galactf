package routes

import (
	"log"

	"github.com/ggmolly/galactf/orm"
	"github.com/ggmolly/galactf/utils"
	"github.com/gofiber/fiber/v2"
)

func GetChallenges(c *fiber.Ctx) error {
	if chals, err := orm.GetChallenges(); err != nil {
		log.Printf("[-] error fetching challenges: %s", err.Error())
		return utils.RestStatusFactoryData(c, fiber.StatusInternalServerError, []orm.Challenge{}, "an error occurred while fetching challenges")
	} else {
		return utils.RestStatusFactoryData(c, fiber.StatusOK, chals, "")
	}
}
