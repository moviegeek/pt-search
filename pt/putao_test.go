package pt

import (
	"os"
	"testing"
)

func TestGetMovieFromSearch(t *testing.T) {
	htmlResult, err := os.Open("putao_search_fixture_test.html")
	defer htmlResult.Close()
	if err != nil {
		t.Fatalf("could not open fixture html file putao_search_fixture_test.html: %v", err)
	}

	movies, err := getMoviesFromSearch(htmlResult)
	if err != nil {
		t.Fatalf("search returns err: %v", err)
	}

	if len(movies) != 4 {
		t.Errorf("expected find 4 movies, but found %d", len(movies))
	}

	for i, title := range []string{
		"[摔跤吧！爸爸] Dangal 2016 BluRay 1080p Atmos TrueHD 7.1 x264-MTeam",
		"[摔跤吧！爸爸] Dangal 2016 BluRay 1080p AVC Atmos TrueHD7 1- Hon3yHD",
		"[摔跤吧！爸爸] Dangal 2016 BluRay 720p AC3 x264-MTeam",
		"[摔跤吧！爸爸] Dangal 2016 Hindi HDTV 1080p x264 DD5 1-Hon3y",
	} {
		if movies[i].Title != title {
			t.Errorf("movie[%d] title does not match\nexpected: %s\nactual:%s", i, title, movies[i].Title)
		}
	}
}
