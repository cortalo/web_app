package infrastructure

import (
	"web_app/domain"

	"github.com/jmoiron/sqlx"
)

type MySQLUserRepository struct {
	db *sqlx.DB
}

func (r *MySQLUserRepository) FindByUsername(username string) (*domain.User, error) {
	sqlStr := `select user_id, username, password from user where username = ?`
	user := domain.User{}
	if err := r.db.Get(&user, sqlStr, username); err != nil {
		return &user, err
	}
	return &user, nil
}

func NewMySQLUserRepository(db *sqlx.DB) *MySQLUserRepository {
	return &MySQLUserRepository{db: db}
}

func (r *MySQLUserRepository) Save(user *domain.User) (err error) {
	sqlStr := `insert into user (user_id, username, password) values (?, ?, ?)`
	_, err = r.db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

func (r *MySQLUserRepository) ExistsByUsername(username string) (bool, error) {
	sqlStr := `select count(user_id) from user where username=?`
	var count int
	if err := r.db.QueryRow(sqlStr, username).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}
