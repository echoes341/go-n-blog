package models

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

type articleDB struct {
	gorm.Model
	Title  string
	Author uint
	Text   string
	Date   time.Time
}

// Article is article struct
type Article struct {
	ID     uint      `json:"id"`
	Title  string    `json:"title"`
	Author uint      `json:"author"`
	Text   string    `json:"text"`
	Date   time.Time `json:"date"`
}

// ArticleGet returns a single article, selected by its ID
func ArticleGet(id int) (*Article, error) {
	var aDb articleDB
	err := db.First(&aDb, id).Error

	article := Article{
		ID:     aDb.ID,
		Title:  aDb.Title,
		Text:   aDb.Text,
		Author: aDb.Author,
		Date:   aDb.Date,
	}

	return &article, err
}

// fillArticle change the articleDB, linked to DB
// to an indipendent struct
func fillArticle(aDb articleDB) Article {
	return Article{
		ID:     aDb.ID,
		Title:  aDb.Title,
		Text:   aDb.Text,
		Author: aDb.Author,
		Date:   aDb.Date,
	}
}

// ArticleCount return a structured count of all the articles
// by year and month. Care: this method is mysql-specific
// For example
// 2017
//   |- 0: 2 (2 articles on january 2017)
func ArticleCount() map[int]map[int]int {
	rows, err := db.Table("article_dbs").Select("YEAR(date) as year, MONTH(date) as month, COUNT(*) as cnt").Group("YEAR(date), MONTH(date)").Rows()
	if err != nil {
		log.Panicln(err)
	}
	result := map[int]map[int]int{}
	for rows.Next() {
		var year, month, count int
		rows.Scan(&year, &month, &count)
		month-- // js date compatibility: months start from 0

		// log.Printf("year: %d  month: %d  count: %d", year, month, count)
		if _, ok := result[year]; !ok { // if map for that year is not initialised
			result[year] = make(map[int]int)
		}
		result[year][month] = count

	}
	// debug
	/*for y, mMap := range result {
		for m, c := range mMap {
			log.Printf("%d %d %d", y, m, c)
		}
	}*/
	return result
}

// Articles returns n articles till date in descending order
func Articles(n int, date time.Time) []Article {
	a := []articleDB{}
	// fetch n first articles in descending order
	db.Where("date <= ?", date).Order("date DESC").Limit(n).Find(&a)
	xa := []Article{}
	for _, aDb := range a {
		article := fillArticle(aDb)
		xa = append(xa, article)
	}

	return xa
}

// ArticleAdd adds an article in database
func ArticleAdd(title, text string, userID uint, date time.Time) Article {
	a := articleDB{
		Author: userID,
		Date:   date,
		Text:   text,
		Title:  title,
	}
	db.NewRecord(a)
	db.Create(&a)

	return fillArticle(a)
}

// ArticleUpdate updates an article and saves it in the database
func ArticleUpdate(id uint, title, text string) Article {
	// begin transaction

	tx := db.Begin()

	var aDb articleDB
	err := tx.First(&aDb, id).Error
	if err != nil {
		log.Printf("[ART-EDIT]: %s\n", err)
		log.Printf("[ART-EDIT]: ID: %d\n", id)
		log.Printf("[ART-EDIT]: Title: %s\n", title)
		log.Printf("[ART-EDIT]: Text: %s\n", text)
		tx.Rollback()
		return Article{}
	}

	aDb.Title = title
	aDb.Text = text

	err = tx.Save(&aDb).Error
	if err != nil {
		log.Printf("[ART-EDIT]: %s\n", err)
		log.Printf("[ART-EDIT]: ID: %d\n", id)
		log.Printf("[ART-EDIT]: Title: %s\n", title)
		log.Printf("[ART-EDIT]: Text: %s\n", text)
		tx.Rollback()
		return Article{}
	}

	// commit transaction
	tx.Commit()
	return fillArticle(aDb)
}

// ArticleRemove removes an Article from the database, selected by ID
func ArticleRemove(id uint) (bool, error) { //bool is for notfound
	// begin transaction
	tx := db.Begin()

	var aDb articleDB
	err := tx.First(&aDb, id).Error
	if err != nil {
		log.Printf("[DEL-ART] Article %d not found. Error: %s", id, err)
		return true, err
	}

	err = tx.Delete(&aDb).Error
	if err != nil {
		log.Printf("[DEL-ART] %s", err)
		return false, err
	}

	tx.Commit()
	return false, err
}
