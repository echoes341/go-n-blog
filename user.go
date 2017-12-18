package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// User is the struct containing the user informations
type User struct {
	ID int
	Name,
	User string
	Password []byte
}

func (u *User) insert() error {
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

// AuthUser takes a username and a password and return a user if login succeeded, otherwhise an error
func AuthUser(username string, password string) (User, error) {
	fmt.Printf("%s, %s", username, password)
	rows, err := db.Query("SELECT ID, NAME, EMAIL, PASSWORD from users WHERE EMAIL=?;", username)
	if err != nil {
		return User{}, fmt.Errorf("Username and/or password do not match")
	}
	u := User{}
	rows.Next()
	rows.Scan(&u.ID, &u.Name, &u.User, &u.Password)
	fmt.Printf("%v, %s", u, password)
	err = bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		return User{}, fmt.Errorf("Username and/or password do not match")
	}
	return u, nil
}

func encryptPassword(u *User) error {
	var e error
	u.Password, e = bcrypt.GenerateFromPassword(u.Password, bcrypt.DefaultCost)
	return e
}
