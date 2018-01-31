package models

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

// IsLiked returns true if an article has been liked by a specific user
func IsLiked(IDArt, IDUser int) bool {
	err := db.First(&likeDB{}, "id_art = ? AND id_user = ?", IDArt, IDUser)
	return err == nil
}

// Likes returns all the likes related to an article
func Likes(IDArt int) (xl []Like, err error) {
	var lsdb []likeDB
	err = db.Find(&lsdb, "id_art = ?", IDArt).Error
	if err != nil {
		return xl, err
	}

	for _, v := range lsdb {
		xl = append(xl, Like{
			ID:     v.ID,
			IDArt:  v.IDArt,
			IDUser: v.IDUser,
			Date:   v.Date,
		})
	}
	return xl, nil
}

// LikesCount counts how many likes has an article
func LikesCount(IDArt uint) (count int) {
	// find and count all the likes of that given article
	db.Where("id_art = ?", IDArt).Find(&[]likeDB{}).Count(&count)
	return count
}
