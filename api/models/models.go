// main models file
// db management

package models

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	// mysql adapter
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

// NewDB loads the new mysql database
// - u is the user
// - p is the password
// - a is the db address
// - d is the db name
func NewDB(u, p, a, d string) error {
	var err error
	par := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", u, p, a, d)
	log.Println(par)
	db, err = gorm.Open("mysql", par)
	// automigrate dbModels
	db.AutoMigrate(&articleDB{})
	db.AutoMigrate(&userDB{})
	db.AutoMigrate(&commentDB{})
	db.AutoMigrate(&likeDB{})
	return err
}

// Close closes the database
func Close() {
	db.Close()
}
