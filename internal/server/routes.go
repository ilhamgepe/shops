package server

import (
	"github.com/ilhamgepe/shops/internal/repositories"
	"github.com/ilhamgepe/shops/internal/server/handlers"
	"github.com/ilhamgepe/shops/internal/services"
)

func (s *FiberServer) RegisterFiberRoutes() {
	publicV1 := s.App.Group("/api/v1")
	protectedV1 := s.App.Group("/api/v1")

	// repositories
	userRepo := repositories.NewUserImpl(s.DB)

	// services
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)

	// register route here
	// auth
	authHandler := handlers.NewAuthHandler(s.Validate, authService)
	authHandler.Innit(publicV1)

	// utils
	userHandler := handlers.NewUserHandler(s.Validate, userService)
	userHandler.Innit(protectedV1)

}
