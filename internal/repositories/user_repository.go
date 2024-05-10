package repositories

import (
	"context"
	"log"

	"github.com/ilhamgepe/shops/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(ctx context.Context, dto models.CreateUserDTO) (user models.User, err error)
	GetUserById(ctx context.Context, id int64) (user models.User, err error)
	GetUserByEmail(ctx context.Context, email string) (user models.User, err error)
}

type user struct {
	db *sqlx.DB
}

func NewUserImpl(db *sqlx.DB) UserRepository {
	return &user{db: db}
}

func (u *user) CreateUser(ctx context.Context, dto models.CreateUserDTO) (user models.User, err error) {
	result, err := u.db.NamedExecContext(ctx, "INSERT INTO users (first_name,last_name,email,password) VALUES (:first_name, :last_name, :email, :password)", dto)
	log.Printf("result namedcontext user %v\n", result)
	if err != nil {
		log.Printf("error namedcontext %v\n", err)
		return
	}

	id, err := result.LastInsertId()
	log.Println(id)
	if err != nil {
		log.Printf("error lasinsetid %v\n", err)
		return
	}

	return u.GetUserById(ctx, id)
}

func (u *user) GetUserById(ctx context.Context, id int64) (user models.User, err error) {
	if err = u.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = ? LIMIT 1", id); err != nil {
		return
	}
	return
}

func (u *user) GetUserByEmail(ctx context.Context, email string) (user models.User, err error) {
	if err = u.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = ? LIMIT 1", email); err != nil {
		return
	}
	return
}
