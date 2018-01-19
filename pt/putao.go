package pt

import (
	"bytes"
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
	"github.com/eliukblau/pixterm/ansimage"
)

const (
	site          = "https://pt.sjtu.edu.cn"
	siteIndex     = "https://pt.sjtu.edu.cn/index.php"
	siteLogin     = "https://pt.sjtu.edu.cn/login.php"
	siteCheckcode = "https://pt.sjtu.edu.cn/getcheckcode.php"
	siteTakeLogin = "https://pt.sjtu.edu.cn/takelogin.php"
	siteTorrent   = "https://pt.sjtu.edu.cn/torrents.php"
)

// Putao represents the torrent search provider for pt.sjtu.edu.cn
type Putao struct {
	username string
	password string
	client   *http.Client
}

// FindAll searchs the query in pt.sjtu.edu.cn and return a list of movies as
// a result
func (pt *Putao) FindAll(query string) ([]Movie, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		log.Printf("search by empty query string, will return all torrents")
	}

	req, _ := http.NewRequest("GET", siteTorrent, nil)
	q := req.URL.Query()
	q.Add("cat401", "1")
	q.Add("cat402", "1")
	q.Add("cat403", "1")
	q.Add("cat406", "1")
	q.Add("cat410", "1")
	q.Add("incldead", "1")
	q.Add("spstate", "0")
	q.Add("inclbookmarked", "0")
	q.Add("picktype", "0")
	q.Add("search_area", "0")
	q.Add("search_mode", "0")
	q.Add("search", query)
	req.URL.RawQuery = q.Encode()

	resp, err := pt.client.Do(req)
	if err != nil {
		log.Printf("failed to send search request %s: %v", req.URL, err)
		return []Movie{}, err
	}

	defer resp.Body.Close()
	movies, err := pt.getMoviesFromSearch(resp.Body)
	if err != nil {
		log.Printf("could not parse movies from search page: %v", err)
		return []Movie{}, err
	}

	return movies, nil
}

func (pt *Putao) getMoviesFromSearch(result io.Reader) ([]Movie, error) {
	movies := []Movie{}
	log.Printf("create root document from response")
	doc, err := goquery.NewDocumentFromReader(result)
	if err != nil {
		log.Printf("failed to create document: %v", err)
		return movies, err
	}

	doc.Find("table.torrents>tbody>tr:nth-child(n+2)").Each(func(i int, s *goquery.Selection) {
		log.Printf("movie %d: %#v", i, s.Text())

		titleElem := s.Find("td:nth-child(2) table.torrentname>tbody>tr>td:first-child>a")
		title, exist := titleElem.Attr("title")
		if !exist {
			title = titleElem.Text()
		}
		idReg := regexp.MustCompile(`id=([0-9]+)&`)
		href, exist := titleElem.Attr("href")
		id := "-1"
		url := ""
		if exist {
			result := idReg.FindStringSubmatch(href)
			if len(result) > 1 || result[1] != "" {
				id = result[1]
			}
			url = fmt.Sprintf("%s/%s", site, href)
		}

		age := s.ChildrenFiltered("td:nth-child(4)").Text()
		size := s.ChildrenFiltered("td:nth-child(5)").Text()
		seeders := s.ChildrenFiltered("td:nth-child(6)").Text()

		movies = append(movies, Movie{From: "putao", ID: id, Title: title, Age: age, Size: size, Seeder: seeders, URL: url})
	})

	return movies, nil
}

func (pt *Putao) tryLogin() error {
	log.Printf("trying to login to %s", siteLogin)

	resp, err := pt.client.Get(site)
	if err != nil {
		log.Printf("failed to send get request to home page %s: %v", site, err)
		return err
	}

	if resp.StatusCode == 200 {
		log.Printf("already logged in, nothing to do")
		return nil
	}

	if resp.StatusCode >= 400 {
		log.Printf("site does not respond correctly, something is wrong")
		return fmt.Errorf("pt site returns %d:%s", resp.StatusCode, resp.Status)
	}

	return pt.login()
}

func (pt *Putao) login() error {
	checkcode, err := pt.promptCheckcode()
	if err != nil {
		log.Printf("could not get checkcode, %v", err)
		return err
	}

	loginData := url.Values{}
	loginData.Add("username", pt.username)
	loginData.Add("password", pt.password)
	loginData.Add("checkcode", checkcode)
	loginDataEncoded := loginData.Encode()

	req, _ := http.NewRequest("POST", siteTakeLogin, bytes.NewBufferString(loginDataEncoded))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(loginDataEncoded)))

	resp, err := pt.client.Do(req)
	if err != nil {
		log.Printf("failed to perform login on %s: %v", siteLogin, err)
		return err
	}

	if resp.StatusCode == 302 && resp.Header.Get("Location") == siteIndex {
		log.Printf("successfully login")
		return nil
	}

	log.Printf("failed to login, either password or checkcode is wrong")
	return fmt.Errorf("failed to login")
}

func (pt *Putao) promptCheckcode() (string, error) {
	resp, err := pt.client.Get(siteCheckcode)
	if err != nil {
		log.Printf("failed to open login page %s: %v", siteCheckcode, err)
		return "", err
	}

	defer resp.Body.Close()

	pix, err := ansimage.NewFromReader(resp.Body, color.RGBA{0, 0, 0, 0}, ansimage.DitheringWithBlocks)

	if err != nil {
		log.Printf("could not create ansimage from response, %v", err)
		return "", nil
	}

	pix.Draw()

	log.Printf("please input answer of above checkcode: ")

	var answer string
	_, err = fmt.Scanln(&answer)
	if err != nil {
		log.Fatal(err)
	}

	return answer, nil
}

func (pt *Putao) dumpCookie() {
	u, _ := url.Parse(site)
	cookies := pt.client.Jar.Cookies(u)

	for i, cookie := range cookies {
		//log.Printf("curCk.Raw=%s", curCk.Raw)
		log.Printf("Cookie [%d]", i)
		log.Printf("Name\t=%s", cookie.Name)
		log.Printf("Value\t=%s", cookie.Value)
		log.Printf("Path\t=%s", cookie.Path)
		log.Printf("Domain\t=%s", cookie.Domain)
		log.Printf("Expires\t=%s", cookie.Expires)
		log.Printf("RawExpires=%s", cookie.RawExpires)
		log.Printf("MaxAge\t=%d", cookie.MaxAge)
		log.Printf("Secure\t=%t", cookie.Secure)
		log.Printf("HttpOnly=%t", cookie.HttpOnly)
		log.Printf("Raw\t=%s", cookie.Raw)
		log.Printf("Unparsed=%s", cookie.Unparsed)
	}
}

// NewPutao creates a new Putao provider, with the given user name and login password
func NewPutao(username, password string, client *http.Client) *Putao {
	putao := &Putao{username, password, client}

	err := putao.tryLogin()
	if err != nil {
		log.Fatal("could not login to pt.sjtu.edu.cn")
	}

	return putao
}
