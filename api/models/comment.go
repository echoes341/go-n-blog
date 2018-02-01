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

func commentDBGet(id uint, tx *gorm.DB) (cDB commentDB, err error) {
	err = tx.First(&cDB, id).Error
	return
}

// CommentsCount returns the number of the comments related to an article
func CommentsCount(IDArt uint) (count int) {
	// count all the records
	db.Where("id_art = ?", IDArt).Find(&[]commentDB{}).Count(&count)
	return count
}

// Comments returns all the comments to an article
func Comments(IDArt int) (xc []Comment, err error) {
	var cDB []commentDB

	err = db.Find(&cDB, "id_art = ?", IDArt).Error
	fmt.Println(cDB)
	if err != nil {
		return xc, err
	}
	for _, v := range cDB {
		xc = append(xc, fillComment(v))
	}
	return xc, nil
}

// CommentAdd adds a new comment in the database
func CommentAdd(aID, uID uint, d time.Time, content string) (c Comment, err error) {
	cDB := commentDB{
		IDArt:   aID,
		IDUser:  uID,
		Date:    d,
		Content: content,
	}
	err = db.Create(&cDB).Error
	if err != nil {
		return
	}
	c = fillComment(cDB)
	return
}

// CommentRemove removes a Comment from the database, selected by ID
func CommentRemove(id uint) (err error) {
	tx := db.Begin()
	var c commentDB
	// look for comment
	c, err = commentDBGet(id, tx)

	if err != nil {
		return err
	}

	err = tx.Delete(&c).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return
}
