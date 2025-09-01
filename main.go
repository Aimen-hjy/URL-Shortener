package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var long_to_short, short_to_long map[string]string

func hash_process_long_to_short(long_url string) string {
	hash := md5.Sum([]byte(long_url))
	res := hex.EncodeToString(hash[:])
	return res[:10]
}
func addUrl(long_url string) {
	short_url, ok := long_to_short[long_url]
	if ok {
		log.Println("short url already exists:", short_url)
	} else {
		short_url = hash_process_long_to_short(long_url)
		long_to_short[long_url] = short_url
		short_to_long[short_url] = long_url
		log.Println("add url:", long_url, "->", short_url)
	}
}
func getLongUrl(short_url string) (string, bool) {
	long_url, ok := short_to_long[short_url]
	return long_url, ok
}
func mainHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}
func PostHandler(c *gin.Context) {
	option := c.PostForm("option")
	url := c.PostForm("url")
	if option == "add" {
		addUrl(url)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"OriginalUrl": url,
			"Result":      long_to_short[url],
		})
	} else if option == "browse" {
		original_url, ok := getLongUrl(url)
		if !ok {
			c.HTML(http.StatusOK, "error.html", gin.H{})
			//TODO: Add a page to show error
		} else {
			c.Redirect(http.StatusFound, original_url)
		}
	}
}
func RedirectHandler(c *gin.Context) {
	short_url := c.Param("short_url")
	original_url, ok := getLongUrl(short_url)
	if !ok {
		c.HTML(http.StatusOK, "error.html", gin.H{})
	} else {
		c.Redirect(http.StatusFound, original_url)
	}
}
func main() {
	short_to_long = make(map[string]string)
	long_to_short = make(map[string]string)
	r := gin.Default()
	r.LoadHTMLGlob("static/*")
	r.GET("/", mainHandler)
	r.POST("/", PostHandler)
	r.GET("/:short_url", RedirectHandler)
	r.Run(":8080")
}
