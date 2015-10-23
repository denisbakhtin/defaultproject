package models

import (
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int64
	Email     string
	Name      string
	Password  []byte
	Timestamp time.Time
}

func (user *User) HashPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = hash
	return nil
}

func GetUser(db *sqlx.DB, id int64) (*User, error) {
	user := User{}
	err := db.Get(&user, "SELECT * FROM users WHERE id=$1", id)
	return &user, err
}

func GetUserByEmail(db *sqlx.DB, email string) (*User, error) {
	user := User{}
	err := db.Get(user, "SELECT * FROM users WHERE lower(email)=lower($1)", email)
	return &user, err
}

func InsertUser(db *sqlx.DB, user *User) error {
	err := db.QueryRow("INSERT INTO users(email, name, password, timestamp) VALUES(lower($1),$2,$3,$4) RETURNING id", user.Email, user.Name, user.Password, time.Now()).Scan(&user.Id)
	return err
}
