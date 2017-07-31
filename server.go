package main

import (
	"log"
	"net/http"
	"os"

	"github.com/juju/persistent-cookiejar"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	"bitbucket.org/laputa/movie-search/pt"
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
	ptPass := getenv("PT_PASS", "")
	if ptPass == "" {
		log.Fatal("please set pt password in PT_PASS environment variable")
	}

	hdcUser := getenv("HDC_USER", "laputa")
	hdcPass := getenv("HDC_PASS", "")
	if hdcPass == "" {
		log.Fatal("please set hdc password in HDC_PASS environment variable")
	}

	cookieFile := getenv("GOCOOKIES", ".cookies")

	cookieJar, _ := cookiejar.New(&cookiejar.Options{Filename: cookieFile})
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: cookieJar,
	}

	putao := pt.NewPutao(ptUser, ptPass, client)
	hdc := pt.NewHDC(hdcUser, hdcPass, client)

	log.Printf("cookie is saved in %s", cookieFile)
	cookieJar.Save()

	log.Printf("started server at %s", port)

	m := martini.Classic()
	m.Use(render.Renderer())

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

		return http.StatusOK, "added to download queue"
	})

	m.RunOnAddr(port)
}
