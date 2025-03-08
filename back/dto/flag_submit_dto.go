package dto

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var (
	dtoValidator   = validator.New()
	ErrInvalidBody = errors.New("invalid request body")
)

type FlagSubmitDTO struct {
	Flag string `body:"flag" validate:"required,min=24,max=48"`
}

func ParseFlagSubmitDTO(c *fiber.Ctx) (*FlagSubmitDTO, error) {
	var dto FlagSubmitDTO
	if err := c.BodyParser(&dto); err != nil {
		return nil, ErrInvalidBody
	}
	if err := dtoValidator.Struct(dto); err != nil {
		return nil, ErrInvalidBody
	}
	return &dto, nil
}
