package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"

	"bitbucket.org/laputa/movie-search/pt"
)

func searchHandler(provider pt.Provider) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		q := req.URL.Query().Get("q")
		movies, err := provider.FindAll(q)
		if err != nil {
			log.Printf("failed to search pt site: %v", err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.NewEncoder(res).Encode(movies)
		if err != nil {
			log.Printf("could not encode movies to json: %v", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)
	}
}

func getenv(env, def string) string {
	value := os.Getenv(env)
	if value == "" {
		return def
	}
	return value
}

func main() {
	ptUser := getenv("PT_USER", "laputa")
	ptPass := getenv("PT_PASS", "")
	if ptPass == "" {
		log.Fatal("please set pt password in PT_PASS environment variable")
	}

	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: cookieJar,
	}

	putao := pt.NewPutao(ptUser, ptPass, client)

	log.Printf("created putao search provider %v", putao)

	log.Printf("started server at :8080")

	http.HandleFunc("/search", searchHandler(putao))

	http.ListenAndServe(":8080", nil)
}
