package factories

import (
	"fmt"
	"strings"

	"github.com/ggmolly/galactf/middlewares"
	"github.com/ggmolly/galactf/orm"
	"github.com/ggmolly/galactf/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/yeqown/go-qrcode/v2"
)

const (
	qrContent = `$FLAG_PLACEHOLDER`
)

func GenerateQuietRiotCode(c *fiber.Ctx) error {
	user := middlewares.ReadUser(c)
	flag := orm.GenerateFlag(user, "quiet riot code")

	c.Set("Content-Type", "text/plain")
	c.Set("Content-Disposition", "attachment; filename=intercepted_data_extraction.txt")
	data := strings.ReplaceAll(qrContent, "$FLAG_PLACEHOLDER", flag)

	qrc, err := qrcode.NewWith(data, qrcode.WithVersion(4))
	if err != nil {
		fmt.Printf("could not generate QRCode: %v", err)
		return utils.RestStatusFactory(c, fiber.StatusInternalServerError, "failed to write response")
	}
	binaryWriter := BinaryWriter{}
	qrc.Save(&binaryWriter)

	if _, err := c.Response().BodyWriter().Write([]byte(binaryWriter.output)); err != nil {
		return utils.RestStatusFactory(c, fiber.StatusInternalServerError, "failed to write response")
	}
	return c.SendStatus(fiber.StatusOK)
}

type BinaryWriter struct {
	output string
}

func (bw *BinaryWriter) Write(mat qrcode.Matrix) error {
	bmp := mat.Bitmap()

	var sb strings.Builder
	sb.Grow(33 * 33)

	for _, row := range bmp {
		for _, val := range row {
			if val {
				sb.WriteString("1")
			} else {
				sb.WriteString("0")
			}
		}
	}
	bw.output = sb.String()
	return nil
}

func (bw *BinaryWriter) Close() error {
	return nil
}
