package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func defineRoutes(router *gin.Engine) {
	v1 := router.Group("/api/gonblog")
	{
		v1.GET("/article/:id", fetchArt)
		v1.GET("/article/:id/like", fetchArtLikes)
		v1.GET("/article/:id/like/isliked/:userid", fetchArtIsLiked)
		v1.GET("/article/:id/comments", fetchArtComments)
	}
}

func fetchArt(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"article": article,
	})
}

func fetchArtLikes(c *gin.Context) {}

func fetchArtComments(c *gin.Context) {}

func fetchArtIsLiked(c *gin.Context) {}
