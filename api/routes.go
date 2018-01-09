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
		"status": http.StatusOK,
		"data":   article,
	})
}

func fetchArtComments(c *gin.Context) {
	IDArt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "ID not valid",
		})
		return
	}

	comments, err := getComments(IDArt)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Comments not found",
		})
		return
	}

	if len(comments) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Comments not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   comments,
	})
}

func fetchArtLikes(c *gin.Context) {
	IDArt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "ID not valid",
		})
		return
	}

	likes, err := getLikes(IDArt)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Likes not found",
		})
		return
	}

	if len(likes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Likes not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   likes,
	})

}

func fetchArtIsLiked(c *gin.Context) {}
