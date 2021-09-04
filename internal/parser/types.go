package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
)

//DigitalFormat blueray, webdl, etc
type DigitalFormat uint8

//DigitalResolution resolution, 4K, 1080p, 720p
type DigitalResolution uint8

//DigitalFileSize file size in Byte
type DigitalFileSize uint64

const (
	//UnknownDigitalFormat unknow foramt
	UnknownDigitalFormat DigitalFormat = iota
	//Blueray blueray disk
	Blueray
	//HDTV recorded from hdtv
	HDTV
	//WebDL web downlaoded source
	WebDL
	//UHDTV recorded from 4K tv
	UHDTV
	//Blueray3D 3D bluray disk
	Blueray3D
)

const (
	//UnknownResolution unknow resolution
	UnknownResolution DigitalResolution = iota
	//FHD 1080p video
	FHD
	//HD 720p video
	HD
	//UHD4K 4K video
	UHD4K
)

var (
	digitalSizeStringMap = map[string]uint64{
		"TiB": 1024 * 1024 * 1024 * 1024,
		"GiB": 1024 * 1024 * 1024,
		"MiB": 1024 * 1024,
		"KiB": 1024,
	}
	digitalSizeSymbols = []string{"TiB", "GiB", "MiB", "KiB"}
)

//MovieInfo metadata for a PT movie item
type MovieInfo struct {
	Title      string
	Year       int
	Group      string
	Source     DigitalFormat
	Resolution DigitalResolution
	Size       DigitalFileSize
}

var digitalFormatToString = map[DigitalFormat]string{
	Blueray:              "Blueray",
	HDTV:                 "HDTV",
	WebDL:                "WebDL",
	UHDTV:                "UHDTV",
	Blueray3D:            "3D",
	UnknownDigitalFormat: "Unknown",
}

var stringToDigitalFormat = map[string]DigitalFormat{
	"Blueray": Blueray,
	"HDTV":    HDTV,
	"WebDL":   WebDL,
	"UHDTV":   UHDTV,
	"3D":      Blueray3D,
	"Unknown": UnknownDigitalFormat,
}

func (s DigitalFileSize) String() string {
	for _, symb := range digitalSizeSymbols {
		if uint64(s) > digitalSizeStringMap[symb] {
			fvalue := float32(s) / float32(digitalSizeStringMap[symb])
			return fmt.Sprintf("%.2f %s", fvalue, symb)
		}
	}
	return fmt.Sprintf("%d Bytes", s)
}

func (f DigitalFormat) String() string {
	return digitalFormatToString[f]
}

// MarshalJSON marshals the enum as a quoted json string
func (f *DigitalFormat) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(digitalFormatToString[*f])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (f *DigitalFormat) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*f = stringToDigitalFormat[j]
	return nil
}

var digitalResolutionToString = map[DigitalResolution]string{
	FHD:               "1080p",
	HD:                "720p",
	UHD4K:             "4K",
	UnknownResolution: "unknown",
}

var stringToDigitalResolution = map[string]DigitalResolution{
	"1080p":   FHD,
	"720p":    HD,
	"4K":      UHD4K,
	"unknown": UnknownResolution,
}

func (f DigitalResolution) String() string {
	return digitalResolutionToString[f]
}

// MarshalJSON marshals the enum as a quoted json string
func (f *DigitalResolution) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(digitalResolutionToString[*f])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (f *DigitalResolution) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*f = stringToDigitalResolution[j]
	return nil
}
