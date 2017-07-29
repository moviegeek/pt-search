package pt

import (
	"os"
	"testing"
)

func TestPutaoGetMovieFromSearch(t *testing.T) {
	htmlResult, err := os.Open("putao_search_fixture_test.html")
	defer htmlResult.Close()
	if err != nil {
		t.Fatalf("could not open fixture html file putao_search_fixture_test.html: %v", err)
	}

	pt := &Putao{}

	movies, err := pt.getMoviesFromSearch(htmlResult)
	if err != nil {
		t.Fatalf("search returns err: %v", err)
	}

	if len(movies) != 4 {
		t.Errorf("expected find 4 movies, but found %d", len(movies))
	}

	for i, title := range []string{} {
		if movies[i].Title != title {
			t.Errorf("movie[%d] title does not match\nexpected: %s\nactual:%s", i, title, movies[i].Title)
		}
	}
	for i, info := range [][]string{
		{
			"[摔跤吧！爸爸] Dangal 2016 BluRay 1080p Atmos TrueHD 7.1 x264-MTeam",
			"16.40GB",
			"1月21天",
			"75",
			"https://pt.sjtu.edu.cn/details.php?id=145014&hit=1&hitfrom=firstpage",
		}, {
			"[摔跤吧！爸爸] Dangal 2016 BluRay 1080p AVC Atmos TrueHD7 1- Hon3yHD",
			"42.75GB",
			"1月21天",
			"1",
			"https://pt.sjtu.edu.cn/details.php?id=145013&hit=1&hitfrom=firstpage",
		}, {
			"[摔跤吧！爸爸] Dangal 2016 BluRay 720p AC3 x264-MTeam",
			"4.37GB",
			"1月21天",
			"168",
			"https://pt.sjtu.edu.cn/details.php?id=145011&hit=1&hitfrom=firstpage",
		}, {
			"[摔跤吧！爸爸] Dangal 2016 Hindi HDTV 1080p x264 DD5 1-Hon3y",
			"9.12GB",
			"2月18天",
			"296",
			"https://pt.sjtu.edu.cn/details.php?id=144410&hit=1&hitfrom=firstpage",
		},
	} {
		if movies[i].Title != info[0] {
			t.Fatalf("movie[%d] title does not match\nexpected: %s\nactual: %s", i, info[0], movies[i].Title)
		}
		if movies[i].Size != info[1] {
			t.Fatalf("movie[%d] size does not match\nexpected: %s\nactual: %s", i, info[1], movies[i].Size)
		}
		if movies[i].Age != info[2] {
			t.Fatalf("movie[%d] age does not match\nexpected: %s\nactual: %s", i, info[2], movies[i].Age)
		}
		if movies[i].Seeder != info[3] {
			t.Fatalf("movie[%d] seeders does not match\nexpected: %s\nactual: %s", i, info[3], movies[i].Seeder)
		}
		if movies[i].URL != info[4] {
			t.Fatalf("movie[%d] url does not match\nexpected: %s\nactual: %s", i, info[4], movies[i].URL)
		}
	}
}

func TestHDCGetMovieFromSearch(t *testing.T) {
	htmlResult, err := os.Open("hdc_search_fixture_test.html")
	defer htmlResult.Close()
	if err != nil {
		t.Fatalf("could not open fixture html file putao_search_fixture_test.html: %v", err)
	}

	hdc := &HDChina{}

	movies, err := hdc.getMoviesFromSearch(htmlResult)
	if err != nil {
		t.Fatalf("search returns err: %v", err)
	}

	if len(movies) != 3 {
		t.Fatalf("expected find 4 movies, but found %d", len(movies))
	}

	for i, info := range [][]string{
		{
			"Dangal.2016.BluRay.1080p.x264.TrueHD.7.1-HDChina",
			"20.11\u00a0GB",
			"1月19天",
			"516",
			"https://hdchina.org/details.php?id=271715&hit=1",
		}, {
			"Dangal.2016.BluRay.720p.x264.DD5.1-HDChina",
			"6.02\u00a0GB",
			"1月19天",
			"310",
			"https://hdchina.org/details.php?id=271706&hit=1",
		}, {
			"Dangal 2016 Hindi 1080p HDTV x264 DD5.1-Hon3y",
			"9.12\u00a0GB",
			"2月18天",
			"387",
			"https://hdchina.org/details.php?id=270373&hit=1",
		},
	} {
		if movies[i].Title != info[0] {
			t.Fatalf("movie[%d] title does not match\nexpected: %s\nactual: %s", i, info[0], movies[i].Title)
		}
		if movies[i].Size != info[1] {
			t.Fatalf("movie[%d] size does not match\nexpected: %s\nactual: %s", i, info[1], movies[i].Size)
		}
		if movies[i].Age != info[2] {
			t.Fatalf("movie[%d] age does not match\nexpected: %s\nactual: %s", i, info[2], movies[i].Age)
		}
		if movies[i].Seeder != info[3] {
			t.Fatalf("movie[%d] seeders does not match\nexpected: %s\nactual: %s", i, info[3], movies[i].Seeder)
		}
		if movies[i].URL != info[4] {
			t.Fatalf("movie[%d] url does not match\nexpected: %s\nactual: %s", i, info[4], movies[i].URL)
		}
	}
}
