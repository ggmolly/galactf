package orm

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func FuzzTestGetUserFromCookie(f *testing.F) {
	corpus := []string{
		"email-token", // A typical case with empty token
		"invalid",     // Invalid token
		"a5a387877edd63d7691fa1adf6101d0d%3A17df361d69156a9a7f2432dd21aafa455239d124196dc873ff38250aeef4c532", // valid token
		"%3A",                   // Edge case with only colon-encoded value
		"short",                 // Another example token
		"another-invalid-token", // Invalid format token
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6IkpvaG4uRG9lQGV4YW1wbGUuY29tIn0", // JWT-like format
	}
	for _, cookie := range corpus {
		f.Add(cookie)
	}

	f.Fuzz(func(t *testing.T, cookie string) {
		app := fiber.New()
		req := app.AcquireCtx(&fasthttp.RequestCtx{})
		req.Cookies("email-token", cookie)
		_, err := GetUserFromCookie(req)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}
