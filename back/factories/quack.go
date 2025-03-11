package factories

import (
	"io"
	"os"

	"github.com/ggmolly/galactf/middlewares"
	"github.com/ggmolly/galactf/orm"
	"github.com/ggmolly/galactf/utils"
	"github.com/gofiber/fiber/v2"
)

const (
	quackOffset = 0x66
	quackLength = 0x20 + 0x06

	quackWriteErr = "failed to write quack image"
)

func GenerateQuackChallenge(c *fiber.Ctx) error {
	user := middlewares.ReadUser(c)

	c.Set("Content-Type", "image/jpeg")
	c.Set("Content-Disposition", "filename=quack.jpg")

	writer := c.Response().BodyWriter()

	file, err := os.OpenFile("./assets/quack.jpg", os.O_RDONLY, 0644)
	if err != nil {
		return utils.RestStatusFactory(c, fiber.StatusInternalServerError, "failed to open quack image")
	}
	defer file.Close()

	// Copy first 0x67 bytes, this is where the "Artist" EXIF tag is stored
	{
		buf := make([]byte, quackOffset)
		if _, err := file.Read(buf); err != nil {
			return utils.RestStatusFactory(c, fiber.StatusInternalServerError, "failed to read quack image")
		}
		if _, err := writer.Write(buf); err != nil {
			return utils.RestStatusFactory(c, fiber.StatusInternalServerError, quackWriteErr)
		}
	}

	// Skip quackLength bytes
	{
		if _, err := file.Seek(quackLength, io.SeekCurrent); err != nil {
			return utils.RestStatusFactory(c, fiber.StatusInternalServerError, "failed to seek")
		}
	}

	flag := orm.GenerateFlag(user, "quack")
	// Copy flag
	{
		if _, err := writer.Write([]byte(flag)); err != nil {
			return utils.RestStatusFactory(c, fiber.StatusInternalServerError, quackWriteErr)
		}
	}

	// Copy the rest of the file
	if _, err := io.Copy(writer, file); err != nil {
		return utils.RestStatusFactory(c, fiber.StatusInternalServerError, quackWriteErr)
	}

	return c.SendStatus(fiber.StatusOK)
}
