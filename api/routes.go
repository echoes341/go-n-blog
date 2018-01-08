package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func fetchArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "ID not valid",
		})
		return
	}

	article, err := getArticle(id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Article not found",
		})
		return
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
