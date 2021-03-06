package models

import (
	"context"

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

// UserMatch takes a user and a password and check in database if
// password match. If so, it returns the selected user (jwt empty),
// otherwise it returns an error and an empty user.
func UserMatch(un, p string) (User, error) {
	uDB := userDB{}
	db.Where("email = ?", un).Or("username = ?", un).Find(&uDB)
	if uDB.ID == 0 {
		return User{}, ErrLoginError // always return an empty user
	}

	// Comparing hash from db with password
	err := bcrypt.CompareHashAndPassword(uDB.PwdHash, []byte(p))
	if err != nil {
		return User{}, ErrLoginError
	}

	// User has been found and password and hash match
	// So we build an empty User model struct and returns it
	return fillUser(uDB), nil
}

type contextKey string

const userContextKey contextKey = "user"

// UserAddToContext adds user to the context
func UserAddToContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userContextKey, u)
}

// UserContext returns the user from the context
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

// UserAdd is for adding a new user inside the database
func UserAdd(username, mail, pwd string) (u User, err error) {
	// check if username is already taken or email is used
	var n int
	err = db.Where("email = ?", mail).Or("username = ?", username).Find(&userDB{}).Count(&n).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	if n != 0 {
		return u, ErrUserPresent
	}
	// calculate hash
	h, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	uDB := userDB{
		Email:    mail,
		Username: username,
		PwdHash:  h,
	}
	db.Create(&uDB)
	return fillUser(uDB), err
}

// UserRemove removes a user from the database
func UserRemove(id uint) (err error) {
	tx := db.Begin()
	var uDB userDB
	err = tx.Where("id = ?", id).Find(&uDB).Error
	if err != nil {
		return
	}

	err = tx.Delete(&uDB).Error
	if err != nil {
		tx.Rollback()
		return
	}
	// does not remove all associated articles and comments
	// maybe it should associate all the orphan models
	// to a dummy user?
	tx.Commit()
	return
}
