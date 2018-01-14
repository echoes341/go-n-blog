package main

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

func getCommentCount(IDArt uint) int {
	i := 0
	// count all the records
	db.Where("id_art = ?", IDArt).Find(&[]commentDB{}).Count(&i)
	return i
}

func getComments(IDArt int) ([]Comment, error) {
	var c []Comment
	var csDB []commentDB

	err := db.Find(&csDB, "id_art = ?", IDArt).Error
	fmt.Println(csDB)
	if err != nil {
		return c, err
	}
	for _, v := range csDB {
		c = append(c, Comment{
			ID:      v.ID,
			IDArt:   v.IDArt,
			IDUser:  v.IDUser,
			Date:    v.Date,
			Content: v.Content,
		})
	}
	return c, nil
}
