package pt

import (
	"bytes"
	"errors"
	"fmt"
	"image/color"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/eliukblau/pixterm/pkg/ansimage"
)

const (
	hdcSite           = "https://hdchina.org"
	hdcSiteIndex      = "https://hdchina.org/index.php"
	hdcSiteLogin      = "https://hdchina.org/login.php"
	hdcSiteImageCheck = "https://hdchina.org/image.php?action=regimage&imagehash=%s"
	hdcSiteTakeLogin  = "https://hdchina.org/takelogin.php"
	hdcSiteTorrent    = "https://hdchina.org/torrents.php"
)

// HDChina represents the torrent search provider for hdchina.org
type HDChina struct {
	username string
	password string
	client   *http.Client
}

// FindAll searchs the query in pt.sjtu.edu.cn and return a list of movies as
// a result
func (hdc *HDChina) FindAll(query string) ([]Movie, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		log.Printf("search by empty query string, will return all torrents")
	}

	req, _ := http.NewRequest("GET", hdcSiteTorrent, nil)
	q := req.URL.Query()
	q.Add("cat17", "1")
	q.Add("cat9", "1")
	q.Add("incldead", "1")
	q.Add("spstate", "0")
	q.Add("inclbookmarked", "0")
	q.Add("boardid", "0")
	q.Add("seeders", "")
	q.Add("search_area", "0")
	q.Add("search_mode", "0")
	q.Add("search", query)
	req.URL.RawQuery = q.Encode()

	resp, err := hdc.client.Do(req)
	if err != nil {
		log.Printf("failed to send search request %s: %v", req.URL, err)
		return []Movie{}, err
	}

	defer resp.Body.Close()
	movies, err := hdc.getMoviesFromSearch(resp.Body)
	if err != nil {
		log.Printf("could not parse movies from search page: %v", err)
		return []Movie{}, err
	}

	return movies, nil
}

func (hdc *HDChina) getMoviesFromSearch(result io.Reader) ([]Movie, error) {
	movies := []Movie{}
	log.Printf("create root document from response")
	doc, err := goquery.NewDocumentFromReader(result)
	if err != nil {
		log.Printf("failed to create document: %v", err)
		return movies, err
	}

	doc.Find("table.torrent_list>tbody>tr:nth-child(n+2)").Each(func(i int, s *goquery.Selection) {
		log.Printf("movie %d: %#v", i, s.Text())

		titleElem := s.Find("td:nth-child(2) table.tbname>tbody>tr>td:nth-child(2)>h3>a")
		title := titleElem.Text()
		idReg := regexp.MustCompile(`id=([0-9]+)&`)
		href, exist := titleElem.Attr("href")
		id := "-1"
		url := ""
		if exist {
			result := idReg.FindStringSubmatch(href)
			if len(result) > 1 || result[1] != "" {
				id = result[1]
			}
			url = fmt.Sprintf("%s/%s", hdcSite, href)
		}

		age := s.ChildrenFiltered("td:nth-child(4)").Text()
		size := s.ChildrenFiltered("td:nth-child(5)").Text()
		seeders := s.ChildrenFiltered("td:nth-child(6)").Text()

		movies = append(movies, Movie{From: "hdc", ID: id, Title: title, Age: age, Size: size, Seeder: seeders, URL: url})
	})

	return movies, nil
}

func (hdc *HDChina) tryLogin() error {
	log.Printf("trying to login to %s", hdcSiteLogin)

	resp, err := hdc.client.Get(hdcSite)
	if err != nil {
		log.Printf("failed to send get request to home page %s: %v", hdcSite, err)
		return err
	}

	if resp.StatusCode == 200 {
		log.Printf("already logged in, nothing to do")
		return nil
	}

	if resp.StatusCode >= 400 {
		log.Printf("hdcSite does not respond correctly, something is wrong")
		return fmt.Errorf("pt hdcSite returns %d:%s", resp.StatusCode, resp.Status)
	}

	return hdc.login()
}

func (hdc *HDChina) login() error {
	imageString, imageHash, err := hdc.promptCheckcode()
	if err != nil {
		log.Printf("could not get imageString, %v", err)
		return err
	}

	loginData := url.Values{}
	loginData.Add("username", hdc.username)
	loginData.Add("password", hdc.password)
	loginData.Add("authcode", "")
	loginData.Add("imagestring", imageString)
	loginData.Add("imagehash", imageHash)
	loginData.Add("trackerssl", "yes")
	loginDataEncoded := loginData.Encode()

	req, _ := http.NewRequest("POST", hdcSiteTakeLogin, bytes.NewBufferString(loginDataEncoded))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(loginDataEncoded)))

	resp, err := hdc.client.Do(req)
	if err != nil {
		log.Printf("failed to perform login on %s: %v", hdcSiteLogin, err)
		return err
	}

	if resp.StatusCode == 302 && resp.Header.Get("Location") == hdcSiteIndex {
		log.Printf("successfully login")
		return nil
	}

	log.Printf("failed to login, either password or checkcode is wrong")
	return fmt.Errorf("failed to login")
}

func (hdc *HDChina) promptCheckcode() (string, string, error) {
	resp, err := hdc.client.Get(hdcSiteLogin)
	if err != nil {
		log.Printf("failed to open login page %s: %v", hdcSiteLogin, err)
		return "", "", err
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Printf("could not create doc from login page response")
		return "", "", err
	}

	imageHashInput := doc.Find("input[name=imagehash]")
	imageHash, exist := imageHashInput.Attr("value")
	if !exist {
		log.Printf("could not find image hash value from login page")
		return "", "", errors.New("failed to find image hash value from login page")
	}
	imageURL := fmt.Sprintf(hdcSiteImageCheck, imageHash)
	log.Printf("get image from %s", imageURL)
	resp, err = hdc.client.Get(imageURL)
	if err != nil {
		log.Printf("failed to get image from hdc %s", imageURL)
		return "", "", err
	}

	defer resp.Body.Close()

	pix, err := ansimage.NewFromReader(resp.Body, color.RGBA{0, 0, 0, 0}, ansimage.DitheringWithBlocks)

	if err != nil {
		log.Printf("could not create ansimage from response, %v", err)
		return "", "", nil
	}

	pix.Draw()

	log.Printf("please input answer of above checkcode: ")

	var answer string
	_, err = fmt.Scanln(&answer)
	if err != nil {
		log.Fatal(err)
	}

	return answer, imageHash, nil
}

// NewHDC creates a new HDChina provider, with the given user name and login password
func NewHDC(username, password string, client *http.Client) *HDChina {
	hdc := &HDChina{username, password, client}

	err := hdc.tryLogin()
	if err != nil {
		log.Fatal("could not login to hdchina.org")
	}

	return hdc
}
