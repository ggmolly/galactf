package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Returns a standard REST API response for responses not returning data
func RestStatusFactory(c *fiber.Ctx, status int, format string, args ...interface{}) error {
	var success bool
	if status >= 200 && status < 300 {
		success = true
	}
	return c.Status(status).JSON(fiber.Map{
		"status":  status,
		"success": success,
		"message": fmt.Sprintf(format, args...),
	})
}

// Returns a standard REST API response for responses returning data
func RestStatusFactoryData(c *fiber.Ctx, status int, data interface{}, format string, args ...interface{}) error {
	var success bool
	if status >= 200 && status < 300 {
		success = true
	}
	return c.Status(status).JSON(fiber.Map{
		"status":  status,
		"success": success,
		"data":    data,
		"message": fmt.Sprintf(format, args...),
	})
}
