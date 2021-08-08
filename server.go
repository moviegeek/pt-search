package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/justlaputa/cookiejar"
	"golang.org/x/net/publicsuffix"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
	"github.com/martini-contrib/render"

	"github.com/moviegeek/pt-search/pt"
)

const (
	// GCLOUD_FUNCTION_URL is the google cloud function for torrent downloader
	GCLOUD_FUNCTION_URL = "https://us-central1-movie-downloader-176704.cloudfunctions.net/movieDownloader"
)

func getenv(env, def string) string {
	value := os.Getenv(env)
	if value == "" {
		return def
	}
	return value
}

func main() {
	port := getenv("PORT", "4000")
	port = ":" + port

	ptUser := getenv("PT_USER", "laputa")
	ptPass := getenv("PT_PASS", "x")
	if ptPass == "" {
		log.Fatal("please set pt password in PT_PASS environment variable")
	}

	hdcUser := getenv("HDC_USER", "laputa")
	hdcPass := getenv("HDC_PASS", "x")
	if hdcPass == "" {
		log.Fatal("please set hdc password in HDC_PASS environment variable")
	}

	jar := cookiejar.NewPersistentJar(&cookiejar.PersistentJarOptions{
		publicsuffix.List, "pt-search-9c319", "Cookies",
	})

	jar.LoadFromFile("cookies.json")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}

	putao := pt.NewPutao(ptUser, ptPass, client)
	hdc := pt.NewHDC(hdcUser, hdcPass, client)

	log.Printf("started server at %s", port)

	m := martini.Classic()

	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"https://storage.googleapis.com"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	m.Use(render.Renderer())

	m.Use(martini.Static("./build", martini.StaticOptions{Fallback: "/index.html", Exclude: "/api/"}))

	m.Get("/api/search", func(req *http.Request, r render.Render) {
		query := req.FormValue("q")
		putaoResult, err := putao.FindAll(query)
		if err != nil {
			log.Printf("failed to search pt site: %v", err)
			r.Error(http.StatusBadRequest)
		}

		hdcResult, err := hdc.FindAll(query)
		if err != nil {
			log.Printf("failed to search hdc site: %v", err)
			r.Error(http.StatusBadRequest)
		}

		r.JSON(http.StatusOK, append(hdcResult, putaoResult...))
	})

	m.Post("/api/queue", func(req *http.Request) (int, string) {
		from := req.FormValue("from")
		if from == "" {
			return http.StatusBadRequest, "you need to specify download from where: pt or hdc"
		}

		id := req.FormValue("id")
		if id == "" {
			return http.StatusBadRequest, "you need to specify the torrent id"
		}

		log.Printf("got download request for torrent %s - %s", from, id)

		err := sendAddTorrentRequest(from, id)
		if err != nil {
			return http.StatusBadRequest, "failed to send add torrent request"
		}

		return http.StatusOK, "added to download queue"
	})

	m.RunOnAddr(port)
}

func sendAddTorrentRequest(from, id string) error {
	client := &http.Client{}

	url := fmt.Sprintf("%s?source=%s&id=%s", GCLOUD_FUNCTION_URL, from, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create new request object %s", err.Error())
		return err
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Failed to send request: %s", err.Error())
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed response: %d %s", resp.StatusCode, resp.Body)
	}

	return nil
}
