package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"os"
)

var UrlMap map[string]string

func httpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}
func process_url(long_url string) string {
	short_url, ok := UrlMap[long_url]
	if ok {
		fmt.Println("Url exists:", short_url)
	} else {

	}
}
func addUrl() {
	for {
		var order, long_url string
		fmt.Scan(&order, &long_url)
		if order == "exit" {
			os.Exit(0)
		}
		if order == "add" {
			short_url := process_url(long_url)
			UrlMap[long_url] = short_url
		}
		if order == "del" {
			delete(UrlMap, long_url)
		}
	}
}

func main() {
	UrlMap = make(map[string]string)
	go addUrl()
	http.HandleFunc("/", httpHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
