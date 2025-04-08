package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func NewRateLimiterMiddleware(max int, expiration time.Duration, errorMsg string) fiber.Handler {
    return limiter.New(limiter.Config{
        Max:        max,
        Expiration: expiration,
        LimitReached: func(c *fiber.Ctx) error {
            return c.SendString(errorMsg)
        },
    })
}