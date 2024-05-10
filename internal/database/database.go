package database

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ilhamgepe/shops/packages/config"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	DB *sqlx.DB
}

func New() *Service {
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.Get.DB_USERNAME, config.Get.DB_PASSWORD, config.Get.DB_HOST, config.Get.DB_PORT, config.Get.DB_DATABASE))
	if err != nil {
		log.Fatal(err)
	}

	db.SetConnMaxLifetime(50 * time.Second)
	db.SetConnMaxIdleTime(20 * time.Second)

	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)

	return &Service{DB: db}
}

func (s *Service) Run() error {
	err := s.DB.Ping()
	if err != nil {
		return err
	}
	log.Println("connected to database")
	return nil
}

func (s *Service) Shutdown() error {
	log.Println("closing database...")
	return s.DB.Close()
}
