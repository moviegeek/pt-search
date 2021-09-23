package main

import (
	"log"
	"net/http"

	"github.com/justlaputa/cookiejar"
	"github.com/moviegeek/pt-search/internal/searcher"
	"golang.org/x/net/publicsuffix"
)

func main() {

	jar := cookiejar.NewPersistentJar(
		&cookiejar.PersistentJarOptions{
			PublicSuffixList:    publicsuffix.List,
			GCPProjectID:        "pt-search-9c319",
			FireStoreCollection: "Cookies",
		},
	)

	jar.LoadFromFile("cookies.json")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}

	hdc := searcher.NewHDC("", "", client)

	movies, err := hdc.Bookmarks()
	if err != nil {
		log.Printf("failed to get bookmarks: %v", err)
		return
	}

	log.Printf("bookmarks: %#v", movies)
}
