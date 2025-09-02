package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

var database *sql.DB

func getLongByShort(short_url string) (string, bool, error) {
	var long_url string
	err := database.QueryRow("SELECT long_url FROM URLmap WHERE short_url=?", short_url).Scan(&long_url)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, err
	}
	return long_url, true, nil
}
func getShortByLong(long_url string) (string, bool, error) {
	var short_url string
	err := database.QueryRow("SELECT short_url FROM URLmap WHERE long_url=?", long_url).Scan(&short_url)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, err
	}
	return short_url, true, nil
}
func insertUrl(long_url, short_url string) error {
	_, err := database.Exec("INSERT INTO URLmap (long_url, short_url) VALUES (?, ?)", long_url, short_url)
	return err
}
func hash_process_long_to_short(long_url string) string {
	hash := md5.Sum([]byte(long_url))
	res := hex.EncodeToString(hash[:])
	return res[:10]
}
func addUrl(long_url string) string {
	short_url, ok, err := getShortByLong(long_url)
	if err != nil {
		log.Fatal(err)
	}
	if ok {
		log.Println("short url already exists:", short_url)
		return short_url
	} else {
		short_url = hash_process_long_to_short(long_url)
		err = insertUrl(long_url, short_url)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("add url:", long_url, "->", short_url)
		return short_url
	}
}

func mainHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}
func PostHandler(c *gin.Context) {
	option := c.PostForm("option")
	url := c.PostForm("url")
	true_url, _, _ := getShortByLong(url)
	if option == "add" {
		true_url = addUrl(url)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"OriginalUrl": url,
			"Result":      true_url,
		})
	} else if option == "browse" {
		original_url, ok, _ := getLongByShort(url)
		if !ok {
			c.HTML(http.StatusOK, "error.html", gin.H{})
		} else {
			c.Redirect(http.StatusFound, original_url)
		}
	}
}
func RedirectHandler(c *gin.Context) {
	short_url := c.Param("short_url")
	original_url, ok, _ := getLongByShort(short_url)
	if !ok {
		c.HTML(http.StatusOK, "error.html", gin.H{})
	} else {
		c.Redirect(http.StatusFound, original_url)
	}
}
func main() {
	var err error
	database, err = sql.Open("sqlite3", "./map.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	_, err = database.Exec("CREATE TABLE IF NOT EXISTS URLmap (long_url TEXT PRIMARY KEY, short_url TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	r.LoadHTMLGlob("static/*")
	r.GET("/", mainHandler)
	r.POST("/", PostHandler)
	r.GET("/:short_url", RedirectHandler)
	r.Run(":8080")
}
