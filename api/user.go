package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type userDB struct {
	gorm.Model
	Username   string
	Email      string
	ProfileURL string
	PwdHash    []byte
}

// User is user structure
type User struct {
	ID         uint   `json:"id,omitempty"`
	Username   string `json:"username"`
	Email      string `json:"email,omitempty"`
	ProfileURL string `json:"profile_url,omitempty"`
	JWT        string `json:"jwt,omitempty"`
}

// match takes a user and a password and check in database if
// password match. If so, it returns the selected user (jwt empty),
// otherwise it returns an error and an empty user.
func match(user, password string) (User, error) {
	uDB := userDB{}
	stdErr := fmt.Errorf("Login error")
	db.Where("username = ?").Find(&uDB)
	if uDB.ID == 0 {
		return User{}, stdErr // always return an empty user
	}

	// Comparing hash from db with password
	err := bcrypt.CompareHashAndPassword(uDB.PwdHash, []byte(password))
	if err != nil {
		return User{}, stdErr
	}

	// User has been found and password and ash match
	// So we build an empty User model struct and returns it

	u := User{
		ID:         uDB.ID,
		Username:   uDB.Username,
		Email:      uDB.Email,
		ProfileURL: uDB.ProfileURL,
		JWT:        "", // empty
	}

	return u, nil

}
