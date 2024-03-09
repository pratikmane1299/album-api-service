package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Sahara", Artist: "DJ Snake", Price: 11.11},
	{ID: "2", Title: "Raja Baja", Artist: "Nucleya", Price: 20.99},
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", addNewAlbum)
	router.GET("/albums/:id", getAlbumById)

	router.Run("localhost:6969")
}

func getAlbums(c *gin.Context) {
	res := gin.H{"success": true, "data": albums}

	c.JSON(http.StatusOK, res)
}

func addNewAlbum(c *gin.Context) {
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	albums = append(albums, newAlbum)
	c.JSON(http.StatusCreated, newAlbum)
}

func getAlbumById(c *gin.Context) {
	id := c.Param("id")

	for _, value := range albums {
		if value.ID == id {
			c.JSON(http.StatusOK, gin.H{"success": true, "data": value})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Album not found"})
}
