package searcher

import (
	"os"
	"testing"
)

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
			"摔跤吧!爸爸/我和我的冠军女儿(台)/摔跤吧!老爸/摔跤家族 *内封台译简繁中字*",
			"20.11\u00a0GB",
			"1月19天",
			"516",
			"https://hdchina.org/details.php?id=271715&hit=1",
		}, {
			"Dangal.2016.BluRay.720p.x264.DD5.1-HDChina",
			"摔跤吧!爸爸/我和我的冠军女儿(台)/摔跤吧!老爸/摔跤家族 *内封台译简繁中字*",
			"6.02\u00a0GB",
			"1月19天",
			"310",
			"https://hdchina.org/details.php?id=271706&hit=1",
		}, {
			"Dangal 2016 Hindi 1080p HDTV x264 DD5.1-Hon3y",
			"摔跤吧！爸爸/我和我的冠军女儿(台) 高清电视版 外挂简繁中字 豆瓣9.2分",
			"9.12\u00a0GB",
			"2月18天",
			"387",
			"https://hdchina.org/details.php?id=270373&hit=1",
		},
	} {
		if movies[i].Title != info[0] {
			t.Fatalf("movie[%d] title does not match\nexpected: %s\nactual: %s", i, info[0], movies[i].Title)
		}
		if movies[i].SubTitle != info[1] {
			t.Fatalf("movie[%d] subtitle does not match\nexpected: %s\nactual: %s", i, info[1], movies[i].SubTitle)
		}
		if movies[i].Size != info[2] {
			t.Fatalf("movie[%d] size does not match\nexpected: %s\nactual: %s", i, info[2], movies[i].Size)
		}
		if movies[i].Age != info[3] {
			t.Fatalf("movie[%d] age does not match\nexpected: %s\nactual: %s", i, info[3], movies[i].Age)
		}
		if movies[i].Seeder != info[4] {
			t.Fatalf("movie[%d] seeders does not match\nexpected: %s\nactual: %s", i, info[4], movies[i].Seeder)
		}
		if movies[i].URL != info[5] {
			t.Fatalf("movie[%d] url does not match\nexpected: %s\nactual: %s", i, info[5], movies[i].URL)
		}
	}
}
