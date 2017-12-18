package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type user struct {
	ID int
	Name,
	User string
	Password []byte
}

func (u *user) insert() error {
	// ID set to NULL so db autoincrements it

	// !!! DOES NOT WORK
	encryptPassword(u)
	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO users (ID, NAME, EMAIL, PASSWORD) VALUES (NULL, '%s', '%s', '%s');", u.Name, u.User, u.Password))
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

	return nil
}

// Login takes a username and a password and return a user if login succeeded, otherwhise an error
func Login(username string, password string) (user, error) {
	u := user{}
	rows, err := db.Query("SELECT * from users WHERE NAME=?;", username)
	if err != nil {
		return u, err
	}

}

func (u user) IsValidPassword(password string) bool {
	return bcrypt.CompareHashAndPassword(u.Password, []byte(password)) == nil
}

func encryptPassword(u *user) error {
	var e error
	u.Password, e = bcrypt.GenerateFromPassword(u.Password, bcrypt.DefaultCost)
	return e
}
