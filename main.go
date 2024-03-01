package main

import (
	"fmt"
	"jwiegmann.de/rest/entity"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "jan.wiegmann"
	password = "foobar"
	dbname   = "albums"
)

var dbClient *sqlx.DB

var albums = []entity.Album{
	{ID: 1, Title: "Blue Train", Artist: "John Coltrane", Price: 56},
	{ID: 2, Title: "Jeru", Artist: "Gerry Mulligan", Price: 17},
	{ID: 3, Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39},
}

func main() {

	initDatabase()
	fillDatabase()

	// Declare routes
	router := gin.Default()

	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	// Start webserver
	router.Run("localhost:8080")
}

/*
* Establish database connection
 */
func initDatabase() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error
	dbClient, err = sqlx.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}
	defer dbClient.Close()

	err = dbClient.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
}

/*
* Fill database with test data
 */
func fillDatabase() {

	for _, album := range albums {

		fmt.Println(album)

		_, err := dbClient.Exec("INSERT INTO album (id, title, artist, price) VALUES ($1, $2, $3, $4)", album.ID, album.Title, album.Artist, album.Price)

		if err != nil {
			return
		}
	}

	fmt.Println("Database filled: ", albums)
}

/*
* Get all albums.
 */
func getAlbums(c *gin.Context) {

	var albums []entity.Album

	dbClient.Select(&albums, "SELECT id, title, artist, price FROM album;")

	c.IndentedJSON(http.StatusOK, albums)
}

/*
* Add a new album.
 */
func postAlbums(c *gin.Context) {

	var newAlbum entity.Album

	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	res, err := dbClient.Exec("INSERT INTO album (id, title, artist, price) VALUES (?, ?, ?, ?);",
		newAlbum.ID, newAlbum.Title, newAlbum.Artist, newAlbum.Price)

	if err != nil {
		return
	}

	id, err := res.LastInsertId()

	if err != nil {
		return
	}

	newAlbum.ID = int(id)

	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumByID(c *gin.Context) {

	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	var album entity.Album

	if err != nil {
		return
	}

	dbClient.Get(&album, "SELECT id, title, artist, price FROM album WHERE id = ?", intId)

	c.IndentedJSON(http.StatusOK, album)
}
