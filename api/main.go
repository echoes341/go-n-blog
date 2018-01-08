package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func init() {
	// open db connection
	var err error
	db, err = gorm.Open("mysql", "gonblog:gonblog@tcp(127.0.0.1:3306)/gonblog?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Panicln("Database connection failed")
	}

	db.AutoMigrate(&articleDB{})
	db.AutoMigrate(&commentDB{})
}

func main() {
	router := gin.Default()

	v1 := router.Group("/api/gonblog")
	{
		v1.GET("/:id", fetchArticle)
	}
	router.Run(":8081")
}

func fetchArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "ID not valid",
		})
		return
	}

	var aDb articleDB
	err = db.First(&aDb, id).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Article not found",
		})
		return
	}

	article := Article{
		ID:     aDb.ID,
		Title:  aDb.Title,
		Text:   aDb.Text,
		Author: aDb.Author,
		Date:   aDb.Date,
	}

	commentsNum, err := getCommentCount(article.ID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Internal server error. Impossible to get comments count.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       http.StatusOK,
		"article":      article,
		"comments_num": commentsNum,
	})
}
