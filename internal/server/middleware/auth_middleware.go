package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ilhamgepe/shops/internal/models"
	"github.com/ilhamgepe/shops/packages/config"
)

func WithAuth(c *fiber.Ctx) error {
	Authorization := c.Get("Authorization", "")
	if Authorization == "" || Authorization == "Bearer " {
		return c.Status(http.StatusUnauthorized).JSON(models.ErrResponse{
			Message: "tolol",
		})
	}
	tokenString := strings.Split(Authorization, " ")

	token, err := jwt.Parse(tokenString[1], func(t *jwt.Token) (interface{}, error) {
		if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println(method)
			return nil, errors.New("invalid method")
		}
		return []byte(config.Get.JWT_SECRET), nil
	})
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusUnauthorized).JSON(models.ErrResponse{
			Message: models.ErrUnauthorized,
		})
	}
	if !token.Valid {
		return c.Status(http.StatusUnauthorized).JSON(models.ErrResponse{
			Message: models.ErrUnauthorized,
		})
	}

	claims := token.Claims.(jwt.MapClaims)
	c.Locals("userId", int64(claims["sub"].(float64)))
	c.Locals("userRole", claims["role"])
	return c.Next()
}

func WithRefresh(c *fiber.Ctx) error {
	Authorization := c.Get("Authorization", "")
	if Authorization == "" || Authorization == "Bearer " {
		return c.Status(http.StatusUnauthorized).JSON(models.ErrResponse{
			Message: "tolol",
		})
	}
	tokenString := strings.Split(Authorization, " ")

	token, err := jwt.Parse(tokenString[1], func(t *jwt.Token) (interface{}, error) {
		if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println(method)
			return nil, errors.New("invalid method")
		}
		return []byte(config.Get.JWT_REFRESH_SECRET), nil
	})
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusUnauthorized).JSON(models.ErrResponse{
			Message: models.ErrUnauthorized,
		})
	}
	if !token.Valid {
		return c.Status(http.StatusUnauthorized).JSON(models.ErrResponse{
			Message: models.ErrUnauthorized,
		})
	}

	claims := token.Claims.(jwt.MapClaims)
	c.Locals("userId", int64(claims["sub"].(float64)))
	c.Locals("userRole", claims["role"])
	return c.Next()
}
