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
	n := 0
	// find and count all the likes of that given article
	err := db.Find(&[]likeDB{}, "id_art = ?", IDArt).Count(&n).Error

	// count all the records
	return n, err
}

func isLiked(IDArt, IDUser int) bool {
	err := db.First(&likeDB{}, "id_art = ? AND id_user = ?", IDArt, IDUser)
	return err == nil
}

func getLikes(IDArt int) ([]Like, error) {
	var l []Like
	var lsdb []likeDB
	err := db.Find(&lsdb, "id_art = ?", IDArt).Error
	if err != nil {
		return l, err
	}

	for _, v := range lsdb {
		l = append(l, Like{
			ID:     v.ID,
			IDArt:  v.IDArt,
			IDUser: v.IDUser,
			Date:   v.Date,
		})
	}
	return l, nil
}
