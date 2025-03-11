package factories

import (
	"io"
	"log"
	"strings"

	"github.com/ggmolly/galactf/middlewares"
	"github.com/ggmolly/galactf/orm"
	"github.com/ggmolly/galactf/utils"
	"github.com/gofiber/fiber/v2"
)

func stringToBrainfuck(input string, writer io.Writer) error {
	prevValue := 0
	for i := 0; i < len(input); i++ {
		currValue := int(input[i])
		diff := currValue - prevValue
		if diff > 0 {
			for diff > 0 {
				if diff >= 8 {
					if _, err := writer.Write([]byte("++++++++")); err != nil {
						return err
					}
					diff -= 8
				} else {
					if _, err := writer.Write([]byte(strings.Repeat("+", diff))); err != nil {
						return err
					}
					diff = 0
				}
			}
		} else {
			for diff < 0 {
				if diff <= -8 {
					if _, err := writer.Write([]byte("--------")); err != nil {
						return err
					}
					diff += 8
				} else {
					if _, err := writer.Write([]byte(strings.Repeat("-", -diff))); err != nil {
						return err
					}
					diff = 0
				}
			}
		}
		if _, err := writer.Write([]byte(".")); err != nil {
			return err
		}
		prevValue = currValue
	}
	return nil
}

func GenerateMoreOrLess(c *fiber.Ctx) error {
	user := middlewares.ReadUser(c)
	flag := orm.GenerateFlag(user, "more or less")
	writer := c.Response().BodyWriter()

	c.Set("Content-Type", "text/plain")
	c.Set("Content-Disposition", "attachment; filename=run_me.bf")

	if err := stringToBrainfuck(flag, writer); err != nil {
		log.Println(err)
		return utils.RestStatusFactory(c, fiber.StatusInternalServerError, "failed to generate brainfuck code")
	}

	return c.SendStatus(fiber.StatusOK)
}