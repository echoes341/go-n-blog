package main

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

type articleDB struct {
	gorm.Model
	Title  string
	Author string
	Text   string
	Date   time.Time
}

// Article is article struct
type Article struct {
	ID     uint      `json:"id"`
	Title  string    `json:"title"`
	Author string    `json:"author"`
	Text   string    `json:"text"`
	Date   time.Time `json:"date"`
}

func getArticle(id int) (*Article, error) {
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

func getArticleCountByYM() map[int]map[int]int {
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

func getArticles(n int, date time.Time) []Article {
	a := []articleDB{}
	// fetch n first articles in descending order
	db.Find(&a, "date <= ?", date).Order("date DESC").Limit(5)
	xa := []Article{}
	for _, aDb := range a {
		article := Article{
			ID:     aDb.ID,
			Title:  aDb.Title,
			Text:   aDb.Text,
			Author: aDb.Author,
			Date:   aDb.Date,
		}
		xa = append(xa, article)
	}

	return xa
}
