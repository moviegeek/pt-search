package parser

import (
	"testing"
)

func TestMinNonNegtive(t *testing.T) {
	cases := []struct {
		input    []int
		expected int
	}{
		{
			[]int{1, 2, 3, 4},
			1,
		},
		{
			[]int{-1, -1, -1},
			-1,
		},
		{
			[]int{-1, 0, 1},
			0,
		},
		{
			[]int{-1, 3, 4},
			3,
		},
		{
			[]int{1},
			1,
		},
		{
			[]int{-1},
			-1,
		},
	}

	for _, c := range cases {
		result := minNonNegtive((c.input)...)
		if result != c.expected {
			t.Errorf("expect min of %v is %d, got %d", c.input, c.expected, result)
		}
	}
}

func TestParseTitle(t *testing.T) {
	titles := []string{
		//From HDC
		"Operation.Red.Sea.2018.1080p.BluRay.x264.DTS-HD.MA.7.1-HDChina",                    // normal 1080p
		"Big.Fish.and.Begonia.2016.BluRay.720p.x264.DTS-HDChina",                            // normal 720p
		"Killer.Joe.2011.Blu-ray.Remux.1080p.AVC.DTS-HD.MA.5.1-OurBit",                      // with Blu-ray
		"Metro.2013.720p.Blu-Ray.x264.DTS.HDCLUB",                                           // without proper group
		"Mudbound.2017.1080p.NF.WEB-DL.DD5.1.x264-NTG",                                      // web-dl
		"Venom 2018 1080p WEB x264 AAC2 0-SHITBOX[11.62 GB]",                                //web-dl with WEB
		"Made in Hong Kong 1997 720p BluRay x264-WiKi",                                      // use space
		"A.One.and.a.Two.2000.720p.BluRay.x264.DTS-zzz@HDC",                                 // with user@group
		"Another WolfCop 2017 BluRay Remux 1080p AVC DTS-HD MA 5.1[4.15 GB]",                // without group
		"Man.in.Black.1997.UHDTV.4K.HEVC-HDCTV[7.33 GB]",                                    // UHDTV
		"The.Longest.Nite.1998.HDTV.1080p.H264.AAC-luobo333[3.5 GB]",                        // HDTV
		"2036.Origin.Unknown.2018.1080p.Blu-ray.AVC.DTS-HD.MA.5.1-Huan@HDSky.iso[21.43 GB]", //Year in Title
		"Arizona 2018 1080p WEB-DL DD5 1 H264-CMRG[2.94 GB]",
		"Thor Ragnarok 2017 3D SBS 720p AVC AC3 5.1[6.84 GB]",
		"Sadako.3D.2.2013.BluRay.1080p.x264.DTS-HD.MA.5.1-HDWinG[8.72 GB]",
		"The.Predator.2018.2160p.WEBRip.HDR.DD5.1.x265-EVO[18.29 GB]", //4K with 2160p
		//From Putao
		"[中间人】Middle Men 2009 1080p BluRay Remux AVC DTS-HD MA 5.1-EPSiLON[29.39 GB]",                         //Chinese name with non-ascii character
		"[除蚤武士] Flea-picking Samurai 2018 1080p BluRay x264 DTS-WiKi[9.86 GB]",                                //normal 1080p
		"[22英里/拳力逃脱(台)/绝地22哩(港)]Mile 22 2018 BluRay 720p x264 DD5 1-HDChina[4.36 GB]",                         //720p
		"[影] Shadow 2018 WEB-DL 1080p H264 AAC-PuTao[2.50 GB]",                                                //webdl
		"[晚熟男人] The Late Bloomer 2016 NF WEBRip 1080p H264 DD 5.1-PuTao[2.51 GB]",                             //webrip
		"[网络谜踪] Searching 2018 1080p Blu-ray DTS-HD MA 5 1 x264-PbK[7.39 GB]",                                 //blu-ray
		"[巨齿鲨] The Meg 2018 1080p BluRay DDP7.1 x264-PuTao",                                                   //without size
		"[协商 / 智命谈判(港) / 极智对决(台) / 谈判]The Negotiation 2018 1080p FHDRip H264 AAC-NonDRM[4.31 GB]",             //fhdrip
		"[碟中谍5：神秘国度] Mission Impossible Rogue Nation 2015 BluRay 1080p AVC Atmos TrueHD7.1-tyx@TTG[33.75 GB]", //@mark
		"[碟中谍6：全面瓦解]Mission Impossible Fallout 2018 BluRay 1080p AVC Atmos TrueHD7 1-MTeam[41.76 GB]",         //without space
		"[2001太空漫游] 2001: A Space Odyssey 1968 720p BluRay DD5.1 x264-Geek[8.71 GB]",                          //with year and comma
		"The Last Chance Diary Of Comedians 2013 JPN Blu-ray 1080p AVC DTS-HD MA 5 1-DiY@KBu[22.76 GB]",       ///without chinese title
		"[铁血战士] The Predator 2018 1080p KORSUB HDRip x264 AAC2 0-STUTTERSHIT[3.54 GB]",                        //hdrip
		"Alpha 2018 WEB-DL 1080p H264 AAC-PuTao",                                                              //without both chinese title and size
		"Godzilla The Planet Eater (2018) [WEBRip] [720p] [YTS AM][774.07 MB]",                                //
		"[龙年] Year of the Dragon 1985 Repack BDRip X264 iNT-TLF[2.54 GB]",
	}

	expected := []MovieInfo{
		//HDC
		{"Operation Red Sea", 2018, "HDChina", Blueray, FHD, 0},
		{"Big Fish and Begonia", 2016, "HDChina", Blueray, HD, 0},
		{"Killer Joe", 2011, "OurBit", Blueray, FHD, 0},
		{"Metro", 2013, "", Blueray, HD, 0},
		{"Mudbound", 2017, "NTG", WebDL, FHD, 0},
		{"Venom", 2018, "SHITBOX", WebDL, FHD, 11620000000},
		{"Made in Hong Kong", 1997, "WiKi", Blueray, HD, 0},
		{"A One and a Two", 2000, "HDC", Blueray, HD, 0},
		{"Another WolfCop", 2017, "", Blueray, FHD, 4150000000},
		{"Man in Black", 1997, "HDCTV", UHDTV, UHD4K, 7330000000},
		{"The Longest Nite", 1998, "luobo333", HDTV, FHD, 3500000000},
		{"2036 Origin Unknown", 2018, "HDSky", Blueray, FHD, 21430000000},
		{"Arizona", 2018, "CMRG", WebDL, FHD, 2940000000},
		{"Thor Ragnarok", 2017, "", Blueray3D, HD, 6840000000},
		{"Sadako", 2013, "HDWinG", Blueray3D, FHD, 8720000000},
		{"The Predator", 2018, "EVO", WebDL, UHD4K, 18290000000},
		//Putao
		{"中间人】Middle Men", 2009, "EPSiLON", Blueray, FHD, 29390000000},
		{"Flea-picking Samurai", 2018, "WiKi", Blueray, FHD, 9860000000},
		{"Mile 22", 2018, "HDChina", Blueray, HD, 4360000000},
		{"Shadow", 2018, "PuTao", WebDL, FHD, 2500000000},
		{"The Late Bloomer", 2016, "PuTao", WebDL, FHD, 2510000000},
		{"Searching", 2018, "PbK", Blueray, FHD, 7390000000},
		{"The Meg", 2018, "PuTao", Blueray, FHD, 0},
		{"The Negotiation", 2018, "NonDRM", UnknownDigitalFormat, FHD, 4310000000},
		{"Mission Impossible Rogue Nation", 2015, "TTG", Blueray, FHD, 33750000000},
		{"Mission Impossible Fallout", 2018, "MTeam", Blueray, FHD, 41760000000},
		{"2001: A Space Odyssey", 1968, "Geek", Blueray, HD, 8710000000},
		{"The Last Chance Diary Of Comedians", 2013, "KBu", Blueray, FHD, 22760000000},
		{"The Predator", 2018, "STUTTERSHIT", UnknownDigitalFormat, FHD, 3540000000},
		{"Alpha", 2018, "PuTao", WebDL, FHD, 0},
		{"Godzilla The Planet Eater", 2018, "", WebDL, HD, 774070000},
		{"Year of the Dragon", 1985, "TLF", Blueray, UnknownResolution, 2540000000},
	}

	for i := range titles {
		result := ParseTitle(titles[i])
		if result != expected[i] {
			t.Errorf("Faied on title %s\nExpected: %+v\nResult  : %+v", titles[i], expected[i], result)
		}
	}
}
