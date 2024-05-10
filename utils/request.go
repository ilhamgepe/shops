package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/ilhamgepe/shops/internal/models"
)

func BodyParser(dto interface{}, ctx *fiber.Ctx) error {
	if err := ctx.BodyParser(&dto); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(models.ErrResponse{
			Errors: models.ErrInternalServerError,
		})
	}
	return nil
}

func ValidationErrorResponse(err error, ctx *fiber.Ctx) error {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		messages := make(map[string]string, len(ve))
		for _, v := range ve {
			switch v.Tag() {
			case "email":
				messages[v.Field()] = fmt.Sprintf("%s is not valid email", v.Value())
			case "required":
				messages[v.Field()] = fmt.Sprintf("%s is required", v.Field())
			case "min":
				messages[v.Field()] = fmt.Sprintf("%s must be at least %s characters", v.Field(), v.Param())
			case "max":
				messages[v.Field()] = fmt.Sprintf("%s must be at most %s characters", v.Field(), v.Param())
			}
		}
		return ctx.Status(http.StatusBadRequest).JSON(models.ErrResponse{
			Errors:  models.ErrBadRequest,
			Message: messages,
			Code:    http.StatusBadRequest,
		})
	}
	return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
		"messages": err.Error(),
	})
}
