package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
)

var long_to_short, short_to_long map[string]string

func hashRedirectHandler(w http.ResponseWriter, r *http.Request) {
	short_url := r.URL.Path
	original_url, ok := short_to_long[short_url[1:]]
	if !ok {
		log.Println("Error:Invalid short url", short_url)
	} else {
		http.Redirect(w, r, original_url, http.StatusFound)
		log.Println("Redirecting to", original_url)
	}
}
func hash_process_long_to_short(long_url string) string {
	hash := md5.Sum([]byte(long_url))
	res := hex.EncodeToString(hash[:])
	return res[:10]
}
func addUrl() {
	for {
		var order, long_url string
		fmt.Scan(&order)
		if order == "exit" {
			os.Exit(0)
		}
		if order == "add" {
			fmt.Scan(&long_url)
			short_url, ok := long_to_short[long_url]
			if ok {
				log.Println("short url already exists:", short_url)
			} else {
				short_url = hash_process_long_to_short(long_url)
				long_to_short[long_url] = short_url
				short_to_long[short_url] = long_url
				log.Println("add url:", long_url, "->", short_url)
			}
			continue
		}
		if order == "del" {
			fmt.Scan(&long_url)
			short_url, ok := long_to_short[long_url]
			if ok {
				delete(long_to_short, long_url)
				delete(short_to_long, short_url)
				log.Println("Delete url:", long_url)
			} else {
				log.Println("long url not exists:", long_url)
			}
			continue
		}
		log.Println("order not exists:", order)
	}
}

func main() {
	short_to_long = make(map[string]string)
	long_to_short = make(map[string]string)
	go addUrl()
	http.HandleFunc("/", hashRedirectHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {

		log.Fatal(err)
	}
}
