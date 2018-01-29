package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type commentDB struct {
	gorm.Model
	IDArt   uint
	IDUser  uint
	Date    time.Time
	Content string
}

// Comment is comment model struct
type Comment struct {
	ID      uint      `json:"id"`
	IDArt   uint      `json:"id_art"`
	IDUser  uint      `json:"id_user"`
	Date    time.Time `json:"date"`
	Content string    `json:"content"`
}

func fillComment(cDB commentDB) Comment {
	return Comment{
		ID:      cDB.ID,
		IDArt:   cDB.IDArt,
		IDUser:  cDB.IDUser,
		Date:    cDB.Date,
		Content: cDB.Content,
	}
}

// CommentsCount returns the number of the comments related to an article
func CommentsCount(IDArt uint) int {
	i := 0
	// count all the records
	db.Where("id_art = ?", IDArt).Find(&[]commentDB{}).Count(&i)
	return i
}

// Comments returns all the comments to an article
func Comments(IDArt int) ([]Comment, error) {
	var c []Comment
	var cDB []commentDB

	err := db.Find(&cDB, "id_art = ?", IDArt).Error
	fmt.Println(cDB)
	if err != nil {
		return c, err
	}
	for _, v := range cDB {
		c = append(c, fillComment(v))
	}
	return c, nil
}
