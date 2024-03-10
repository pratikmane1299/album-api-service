package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
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

var db *sql.DB

func main() {

	cfg := mysql.Config{
		User:   "root",
		Passwd: "root",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}

	var err error

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal("Could not connect db")
	}

	fmt.Println("Connected to db :)")

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", addNewAlbum)
	router.GET("/albums/:id", getAlbumById)
	router.PATCH("/albums/:id", updateAlbumById)
	router.DELETE("/albums/:id", deleteAlbumById)

	router.Run("localhost:6969")
}

func getAlbums(c *gin.Context) {
	var albums []album

	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		fmt.Printf("getAlbums - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "something went wrong"})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var alb album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			fmt.Printf("getAlbums - %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "something went wrong"})
			return
		}

		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("getAlbums - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "something went wrong"})
		return
	}

	res := gin.H{"success": true, "data": albums}

	c.JSON(http.StatusOK, res)
}

func addNewAlbum(c *gin.Context) {
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		fmt.Printf("Error: addNewAlbum - %v", err)
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "something went wrong"})
		return
	}

	result, err := db.Exec("INSERT INTO album(title, artist, price) VALUES(?, ?, ?)", &newAlbum.Title, &newAlbum.Artist, &newAlbum.Price)
	if err != nil {
		fmt.Printf("Error: addNewAlbum - %v", err)
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "something went wrong"})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("Error: addNewAlbum - %v", err)
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "something went wrong"})
		return
	}

	albums = append(albums, newAlbum)
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": id, "message": "new album created"})
}

func getAlbumById(c *gin.Context) {
	id := c.Param("id")

	var alb album

	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Error: getAlbumById - %v", err)
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "album not found"})
			return
		}

		fmt.Printf("Error: getAlbumById - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": alb})
}

func updateAlbumById(c *gin.Context) {
	id := c.Param("id")
	var alb album

	if err := c.BindJSON(&alb); err != nil {
		fmt.Printf("Error: updateAlbumById - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "something went wrong"})
		return
	}

	var albumToUpdate album
	row := db.QueryRow("SELECT id FROM album WHERE id = ?", id)
	if err := row.Scan(&albumToUpdate.ID); err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Error: getAlbumById - %v", err)
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "album not found"})
			return
		}

		fmt.Printf("Error: updateAlbumById - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "something went wrong"})
		return
	}

	result, err := db.Exec("UPDATE Album SET title = ?, artist = ?, price = ? WHERE id = ?", alb.Title, alb.Artist, alb.Price, id)
	if err != nil {
		fmt.Printf("Error: updateAlbumById - %v", err)
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "something went wrong"})
		return
	}

	totalUpdatedRows, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Error: updateAlbumById - %v", err)
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"succes": true, "data": totalUpdatedRows, "message": "album updated"})
}

func deleteAlbumById(c *gin.Context) {
	id := c.Param("id")

	var album album

	row := db.QueryRow("SELECT id FROM album WHERE id = ?", id)
	if err := row.Scan(&album.ID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "album not found"})
		return
	}

	result, err := db.Exec("DELETE FROM album WHERE id = ?", id)
	if err != nil {
		fmt.Printf("Error: deleteAlbumById - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "something went wrong"})
		return
	}

	deletedRows, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Error: deleteAlbumById - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": deletedRows, "message": "album deleted"})
}
