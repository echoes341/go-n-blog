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
	for rows.Next() {
		var year, month, count int
		rows.Scan(&year, &month, &count)
		log.Printf("year: %d  month: %d  count: %d", year, month, count)
	}

	result := map[int]map[int]int{}
	result[2017] = make(map[int]int)
	result[2017][9] = 2
	_, ok := result[2016]
	log.Printf("%v", ok)
	for y, mMap := range result {
		for m, c := range mMap {
			log.Printf("%d %d %d", y, m, c)
		}
	}
	return result
}
