package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

type likeDB struct {
	gorm.Model
	IDArt  uint
	IDUser uint
	Date   time.Time
}

// Like is a like model struct
type Like struct {
	ID     uint      `json:"id"`
	IDArt  uint      `json:"id_art"`
	IDUser uint      `json:"id_user"`
	Date   time.Time `json:"date"`
}

func getLikeCount(IDArt int) (int, error) {
	i := 0
	// find and count all the likes of that given article
	err := db.Find(&[]likeDB{}, "id_art = ?", IDArt).Count(&i).Error

	// count all the records
	return i, err
}

func isLiked(IDArt, IDUser int) bool {
	err := db.First(&likeDB{}, "id_art = ? AND id_user = ?", IDArt, IDUser)
	return err == nil
}
