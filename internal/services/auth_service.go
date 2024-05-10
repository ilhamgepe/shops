package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ilhamgepe/shops/internal/models"
	"github.com/ilhamgepe/shops/internal/repositories"
	"github.com/ilhamgepe/shops/packages/config"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	ur repositories.UserRepository
}

type JwtClaims struct {
	jwt.RegisteredClaims
	Sub  int64       `json:"sub"`
	Role models.Role `json:"role"`
}

func NewAuthService(ur repositories.UserRepository) *AuthService {
	return &AuthService{ur: ur}
}

func (as *AuthService) Register(ctx context.Context, dto models.CreateUserDTO) (models.User, error) {
	return as.ur.CreateUser(ctx, dto)
}

func (as *AuthService) Login(ctx context.Context, dto models.LoginUserDTO) (models.User, error) {
	return as.ur.GetUserByEmail(ctx, dto.Email)
}

func (as *AuthService) Me(ctx context.Context, id int64) (models.User, error) {
	return as.ur.GetUserById(ctx, id)
}

// helpers

func (as *AuthService) CreateToken(id int64, role models.Role, exp time.Duration, secret []byte) (string, error) {
	claims := JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Get.APP_NAME,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Sub:  id,
		Role: role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func (as *AuthService) ValidateToken(tokenString string, secret []byte) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid method")
		}
		return secret, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}

func (as *AuthService) GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (as *AuthService) ComparePassword(plain string, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err == nil
}
