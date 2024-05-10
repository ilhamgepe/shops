package services

import (
	"log"
	"testing"
	"time"

	"github.com/ilhamgepe/shops/internal/database"
	"github.com/ilhamgepe/shops/internal/models"
	"github.com/ilhamgepe/shops/internal/repositories"
	"github.com/ilhamgepe/shops/packages/config"
	"github.com/ilhamgepe/shops/packages/logs"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

var (
	db *sqlx.DB
	ur repositories.UserRepository
	as *AuthService
)

func init() {
	config.Load("../../.env")
	logs.NewLogger("../../logs_test.log")
	mysqlConn := database.New()
	db = mysqlConn.DB
	ur = repositories.NewUserImpl(db)
	as = NewAuthService(ur)
}

func TestToken(t *testing.T) {
	mockUser := models.User{
		ID:        1,
		FirstName: "ilham",
		LastName:  nil,
		Email:     "ilham@gmail.com",
		Password:  "kzmaacks98",
		Role:      models.Role("user"),
	}
	var token string
	var err error
	t.Run("create token", func(t *testing.T) {
		token, err = as.CreateToken(mockUser.ID, mockUser.Role, 1*time.Hour, []byte(config.Get.JWT_SECRET))
		require.Nil(t, err)
		require.NotNil(t, token)
	})

	t.Run("validate token", func(t *testing.T) {
		err := as.ValidateToken(token, []byte(config.Get.JWT_SECRET))
		log.Println(err)
		require.Nil(t, err)
	})
}

func TestPassword(t *testing.T) {
	var hashedPass string
	var err error
	t.Run("create password", func(t *testing.T) {
		hashedPass, err = as.GenerateHashPassword("hello_world")
		require.Nil(t, err)
	})
	t.Run("compare password", func(t *testing.T) {
		isMatch := as.ComparePassword("hello_world", hashedPass)
		require.Equal(t, true, isMatch)
	})
}
