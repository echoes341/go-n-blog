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

func fetchLike(a, u uint, database *gorm.DB) (lDB likeDB, err error) {
	err = database.First(&lDB, "id_art = ? AND id_user = ?", a, u).Error
	return lDB, err
}

// IsLiked returns true if an article has been liked by a specific user
func IsLiked(aID, uID uint) bool {
	_, err := fetchLike(aID, uID, db)
	return err != nil
}

// LikeToggle set a like by a user or removes it
func LikeToggle(aID, uID uint) (added bool, err error) {
	tx := db.Begin()

	// Check if article exist, otherwise return error
	_, err = articleGet(aID, tx)
	if err != nil {
		// Rollback is not necessary as we have not edited the database
		return
	}

	var lDB likeDB
	// search if liked is present
	if lDB, err = fetchLike(aID, uID, tx); err == nil {
		err = tx.Delete(&lDB).Error
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
		return
	}

	lDB = likeDB{
		IDArt:  aID,
		IDUser: uID,
		Date:   time.Now(),
	}
	if new := tx.NewRecord(lDB); !new {
		return
	}
	err = tx.Create(&lDB).Error
	if err != nil {
		tx.Rollback()
		return
	}
	added = true
	tx.Commit()
	return
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
