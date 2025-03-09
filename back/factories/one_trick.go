package factories

import "github.com/gofiber/fiber/v2"

func RenderOneTrick(c *fiber.Ctx) error {
	return c.Render("one_trick/index", fiber.Map{})
}
