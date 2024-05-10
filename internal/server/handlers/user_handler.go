package handlers

import (
	"database/sql"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/ilhamgepe/shops/internal/models"
	"github.com/ilhamgepe/shops/internal/server/middleware"
	"github.com/ilhamgepe/shops/internal/services"
)

type UserHandler struct {
	*validator.Validate
	*services.UserService
}

func NewUserHandler(v *validator.Validate, s *services.UserService) *UserHandler {
	return &UserHandler{Validate: v, UserService: s}
}

func (uh *UserHandler) Innit(r fiber.Router) {
	r.Get("/me", middleware.WithAuth, uh.MeHandler)
}

func (uh *UserHandler) MeHandler(c *fiber.Ctx) error {
	user, err := uh.UserService.GetById(c.Context(), c.Locals("userId").(int64))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(http.StatusBadRequest).JSON(models.ErrResponse{
				Errors:  models.ErrBadRequest,
				Message: models.ErrUserNotRegistered,
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(models.ErrResponse{
			Errors: models.ErrInternalServerError,
		})
	}
	return c.Status(http.StatusOK).JSON(models.SuccessResponse{
		Data: user,
	})
}
