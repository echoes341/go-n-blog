package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

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
	defer db.Close()

	db.AutoMigrate(&articleDB{})
}

type articleDB struct {
	gorm.Model
	Title  string    `json:"title"`
	Author string    `json:"author"`
	Text   string    `json:"text"`
	Date   time.Time `json:"date"`
}

// Article is article struct
type Article struct {
	ID     int       `json:"id"`
	Title  string    `json:"title"`
	Author string    `json:"author"`
	Text   string    `json:"text"`
	Date   time.Time `json:"date"`
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

	article := articleDB{}
	db.First(&article)
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"article": article,
		"id":      id,
	})
}
