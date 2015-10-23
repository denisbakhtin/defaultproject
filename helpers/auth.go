package helpers

import (
	"github.com/denisbakhtin/defaultproject/models"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func Login(db *sqlx.DB, email string, password string) (*models.User, error) {
	user := models.User{}
	err := db.Get(&user, "SELECT * FROM users where email=lower($1)", email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return nil, err
	}
	return &user, nil
}
