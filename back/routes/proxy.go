package routes

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ggmolly/galactf/middlewares"
	"github.com/ggmolly/galactf/orm"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

type ProxiedChallengeSettings struct {
	Hostname string
	Port     int
}

var (
	proxySettings = map[string]map[string]ProxiedChallengeSettings{
		"bobby_library": {
			"dev":  {Hostname: "127.0.0.1", Port: 10000},
			"prod": {Hostname: "bobby_library", Port: 8080},
		},
	}
)

func ProxyFactory(containerName, baseURL, prettyChalName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := middlewares.ReadUser(c)
		flag := orm.GenerateFlag(user, prettyChalName)

		conf := proxySettings[containerName][os.Getenv("MODE")]
		remotePath := fmt.Sprintf(
			"http://%s:%d/%s",
			conf.Hostname,
			conf.Port,
			c.OriginalURL()[18+len(baseURL)-1:],
		)

		// Inject flag & user ID in the request
		c.Request().Header.Set("X-User-ID", fmt.Sprint(user.ID))
		c.Request().Header.Set("X-GalaCTF-Flag", flag)

		// Just forward the request to the corresponding container
		return proxy.DoTimeout(c, remotePath, time.Second*5)
	}
}
