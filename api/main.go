package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func init() {
	// open db connection
	var err error
	db, err = gorm.Open("mysql", "gonblog:gonblog@tcp(127.0.0.1:3306)/gonblog")
	if err != nil {
		log.Panicln("Database connection failed")
	}

}

func main() {
	router := gin.Default()

	v1 := router.Group("/api/gonblog")
	{
		v1.GET("/:id", fetchArticle)
	}
	router.Run()
}

func fetchArticle(c *gin.Context) {

}
