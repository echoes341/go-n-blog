package main

import (
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

func getCommentCount(IDArt uint) (int, error) {
	i := 0
	// comments := []commentDB{}
	// find all the comments of that givent article
	err := db.Find(&[]commentDB{}, "id_art = ?", IDArt).Count(&i).Error

	// count all the records
	return i, err
}
