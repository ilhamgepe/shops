package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/ilhamgepe/shops/internal/models"
	"github.com/ilhamgepe/shops/internal/server/middleware"
	"github.com/ilhamgepe/shops/internal/services"
	"github.com/ilhamgepe/shops/packages/config"
	"github.com/ilhamgepe/shops/packages/logs"
	"github.com/ilhamgepe/shops/utils"
)

type AuthHandler struct {
	*validator.Validate
	*services.AuthService
}

func NewAuthHandler(v *validator.Validate, s *services.AuthService) *AuthHandler {
	return &AuthHandler{Validate: v, AuthService: s}
}

func (ah *AuthHandler) Innit(r fiber.Router) {
	r.Post("/auth/register", ah.RegisterHandler)
	r.Post("/auth/login", ah.LoginHandler)
	r.Get("/auth/refresh", middleware.WithRefresh, ah.RefreshHandler)
}

func (ah *AuthHandler) RegisterHandler(c *fiber.Ctx) error {
	/* validate request */
	var dto models.CreateUserDTO
	utils.BodyParser(&dto, c)
	if err := ah.Validate.Struct(&dto); err != nil {
		return utils.ValidationErrorResponse(err, c)
	}
	pass, err := ah.AuthService.GenerateHashPassword(dto.Password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.ErrResponse{
			Errors: models.ErrInternalServerError,
		})
	}
	dto.Password = pass
	/* end validate request */

	/* create user */
	user, err := ah.AuthService.Register(c.Context(), dto)
	if err != nil {
		if mErr, ok := err.(*mysql.MySQLError); ok {
			switch mErr.Number {
			case 1062:
				return c.Status(http.StatusBadRequest).JSON(models.ErrResponse{
					Errors:  models.ErrBadRequest,
					Message: "User Already exist",
				})
			default:
				logs.Logger.Error().Stack().Err(errors.New(err.Error())).Msg("failed to register")
				return c.Status(http.StatusInternalServerError).JSON(models.ErrResponse{
					Errors: models.ErrInternalServerError,
				})
			}
		}
		return c.Status(http.StatusInternalServerError).JSON(models.ErrResponse{
			Errors: models.ErrInternalServerError,
		})
	}
	/* end create user */

	/* create token */
	token, err := ah.AuthService.CreateToken(user.ID, user.Role, time.Duration(config.Get.JWT_EXPIRY)*time.Hour, []byte(config.Get.JWT_SECRET))
	if err != nil {
		logs.Logger.Error().Stack().Err(errors.New(err.Error())).Msg("FAILED CREATE TOKEN IN REGISTER")
		return c.Status(http.StatusInternalServerError).JSON(models.ErrResponse{
			Errors: models.ErrInternalServerError,
		})
	}
	refresh, err := ah.AuthService.CreateToken(user.ID, user.Role, time.Duration(config.Get.JWT_REFRESH_EXPIRY)*time.Hour, []byte(config.Get.JWT_REFRESH_SECRET))
	if err != nil {
		logs.Logger.Error().Stack().Err(errors.New(err.Error())).Msg("FAILED CREATE REFRESH TOKEN IN REGISTER")
		return c.Status(http.StatusInternalServerError).JSON(models.ErrResponse{
			Errors: models.ErrInternalServerError,
		})
	}
	/* end create token */

	return c.Status(http.StatusOK).JSON(models.SuccessResponse{
		Data: map[string]any{
			"user":          user,
			"token":         token,
			"refresh_token": refresh,
		},
	})
}
func (ah *AuthHandler) LoginHandler(c *fiber.Ctx) error {
	/* validate request */
	var dto models.LoginUserDTO
	utils.BodyParser(&dto, c)
	if err := ah.Validate.Struct(&dto); err != nil {
		return utils.ValidationErrorResponse(err, c)
	}
	/* end validate request */

	/* find user */
	user, err := ah.AuthService.Login(c.Context(), dto)
	if err != nil {
		log.Println(err)
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
	if match := ah.AuthService.ComparePassword(dto.Password, user.Password); !match {
		return c.Status(http.StatusUnauthorized).JSON(models.ErrResponse{
			Errors:  models.ErrUnauthorized,
			Message: models.ErrInvalidCredentials,
		})
	}
	/* end find user */

	/* create token */
	token, err := ah.AuthService.CreateToken(user.ID, user.Role, time.Duration(config.Get.JWT_EXPIRY)*time.Hour, []byte(config.Get.JWT_SECRET))
	if err != nil {
		logs.Logger.Error().Stack().Err(errors.New(err.Error())).Msg("FAILED CREATE TOKEN IN REGISTER")
		return c.Status(http.StatusInternalServerError).JSON(models.ErrResponse{
			Errors: models.ErrInternalServerError,
		})
	}
	refresh, err := ah.AuthService.CreateToken(user.ID, user.Role, time.Duration(config.Get.JWT_REFRESH_EXPIRY)*time.Hour, []byte(config.Get.JWT_REFRESH_SECRET))
	if err != nil {
		logs.Logger.Error().Stack().Err(errors.New(err.Error())).Msg("FAILED CREATE REFRESH TOKEN IN REGISTER")
		return c.Status(http.StatusInternalServerError).JSON(models.ErrResponse{
			Errors: models.ErrInternalServerError,
		})
	}
	/* end create token */

	return c.Status(http.StatusOK).JSON(models.SuccessResponse{
		Data: map[string]any{
			"user":          user,
			"token":         token,
			"refresh_token": refresh,
		},
	})
}

func (ah *AuthHandler) RefreshHandler(c *fiber.Ctx) error {
	id := c.Locals("userId").(int64)
	user, err := ah.AuthService.Me(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.ErrResponse{
			Errors: models.ErrInternalServerError,
		})
	}

	/* create token */
	token, err := ah.AuthService.CreateToken(user.ID, user.Role, time.Duration(config.Get.JWT_EXPIRY)*time.Hour, []byte(config.Get.JWT_SECRET))
	if err != nil {
		logs.Logger.Error().Stack().Err(errors.New(err.Error())).Msg("FAILED CREATE TOKEN IN REGISTER")
		return c.Status(http.StatusInternalServerError).JSON(models.ErrResponse{
			Errors: models.ErrInternalServerError,
		})
	}
	refresh, err := ah.AuthService.CreateToken(user.ID, user.Role, time.Duration(config.Get.JWT_REFRESH_EXPIRY)*time.Hour, []byte(config.Get.JWT_REFRESH_SECRET))
	if err != nil {
		logs.Logger.Error().Stack().Err(errors.New(err.Error())).Msg("FAILED CREATE REFRESH TOKEN IN REGISTER")
		return c.Status(http.StatusInternalServerError).JSON(models.ErrResponse{
			Errors: models.ErrInternalServerError,
		})
	}
	/* end create token */

	return c.Status(http.StatusOK).JSON(models.SuccessResponse{
		Data: map[string]any{
			"user":          user,
			"token":         token,
			"refresh_token": refresh,
		},
	})
}
