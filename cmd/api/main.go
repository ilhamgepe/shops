package main

import (
	"context"
	"log"
	"os"

	"github.com/ilhamgepe/shops/internal/database"
	"github.com/ilhamgepe/shops/internal/server"
	"github.com/ilhamgepe/shops/packages/config"
	"github.com/ilhamgepe/shops/packages/logs"
	"github.com/ilhamgepe/shops/utils"
)

func init() {
	config.Load(".env")
	logs.NewLogger("logs.log")
}

func main() {
	mysqlConn := database.New()
	mysqlConn.Run()

	srv := server.New(mysqlConn.DB)
	srv.RegisterFiberRoutes()

	go func() {
		if err := srv.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	wait := utils.GracefulShutdown(context.Background(), map[string]utils.Operation{
		"server": func(ctx context.Context) error {
			return srv.Shutdown()
		},
		"database": func(ctx context.Context) error {
			return mysqlConn.Shutdown()
		},
		// "logs": func(ctx context.Context) error {
		// 	return logger.CloseLogger()
		// },
	})
	<-wait

	log.Println("gracefully shutdown")
	os.Exit(0)
}
