package routes

import (
	"fmt"
	"os"
	"strings"
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
		"unsecure_notes": {
			"dev":  {Hostname: "127.0.0.1", Port: 10001},
			"prod": {Hostname: "unsecure_notes", Port: 8080},
		},
		"just_a_moment": {
			"dev":  {Hostname: "127.0.0.1", Port: 10002},
			"prod": {Hostname: "just_a_moment", Port: 8080},
		},
		"cookie_monster": {
			"dev":  {Hostname: "127.0.0.1", Port: 10003},
			"prod": {Hostname: "cookie_monster", Port: 8080},
		},
		"cookie_monster_squared": {
			"dev":  {Hostname: "127.0.0.1", Port: 10004},
			"prod": {Hostname: "cookie_monster_squared", Port: 8080},
		},
	}
)

func ProxyFactory(containerName, baseURL, prettyChalName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := middlewares.ReadUser(c)
		flag := orm.GenerateFlag(user, prettyChalName)
		conf := proxySettings[containerName][os.Getenv("MODE")]
		remotePath := fmt.Sprintf(
			"%s:%d/%s",
			conf.Hostname,
			conf.Port,
			c.OriginalURL()[18+len(baseURL)-1:],
		)

		// Hacky fix for double leading slashes in some cases
		remotePath = "http://" + strings.ReplaceAll(remotePath, "//", "/")

		// Inject flag & user ID in the request
		c.Request().Header.Set("X-User-ID", fmt.Sprint(user.ID))
		c.Request().Header.Set("X-GalaCTF-Flag", flag)
		c.Request().Header.Set("X-Root-Uri", "/api/v1/factories/"+containerName)

		fmt.Println("proxying : ", containerName)

		// Just forward the request to the corresponding container
		return proxy.DoTimeout(c, remotePath, time.Second*5)
	}
}
