package main

import (
	"context"
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
	IsAdmin    bool
}

// User is user structure
type User struct {
	ID         uint   `json:"id,omitempty"`
	Username   string `json:"username"`
	Email      string `json:"email,omitempty"`
	ProfileURL string `json:"profile_url,omitempty"`
	IsAdmin    bool   `json:"is_admin"`
}

func fillUser(uDB userDB) User {
	return User{
		ID:         uDB.ID,
		Username:   uDB.Username,
		Email:      uDB.Email,
		ProfileURL: uDB.ProfileURL,
		IsAdmin:    uDB.IsAdmin,
	}
}

// match takes a user and a password and check in database if
// password match. If so, it returns the selected user (jwt empty),
// otherwise it returns an error and an empty user.
func UserMatch(un, p string) (User, error) {
	uDB := userDB{}
	stdErr := fmt.Errorf("Login error")
	db.Where("email = ?", un).Or("username = ?", un).Find(&uDB)
	if uDB.ID == 0 {
		return User{}, stdErr // always return an empty user
	}

	// Comparing hash from db with password
	err := bcrypt.CompareHashAndPassword(uDB.PwdHash, []byte(p))
	if err != nil {
		return User{}, stdErr
	}

	// User has been found and password and ash match
	// So we build an empty User model struct and returns it

	return fillUser(uDB), nil

}

type contextKey string

const userContextKey contextKey = "user"

func UserAddToContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userContextKey, u)
}

func UserContext(ctx context.Context) *User {
	/*if u, ok := ctx.Value(userContextKey).(*User); ok {
		return u
	}
	return &User{}*/
	return &User{
		ID:       2,
		Username: "debug",
		IsAdmin:  true,
	}
}
