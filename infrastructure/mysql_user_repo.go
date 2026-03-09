package infrastructure

import (
	"web_app/domain"

	"github.com/jmoiron/sqlx"
)

type MySQLUserRepository struct {
	db *sqlx.DB
}

func NewMySQLUserRepository(db *sqlx.DB) *MySQLUserRepository {
	return &MySQLUserRepository{db: db}
}

func (r *MySQLUserRepository) Save(user domain.User) error {
	return nil
}

func (r *MySQLUserRepository) FindByID(id string) (domain.User, error) {
	return domain.User{}, nil
}

func (r *MySQLUserRepository) ExistsByEmail(email string) (bool, error) {
	return false, nil
}
