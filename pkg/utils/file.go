package utils

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"strings"
)

type File string

// Filename returns name of file
func (fp File) Filename() string {
	if fp == "" {
		return ""
	}

	file := strings.Split(string(fp), "/")
	last := file[len(file)-1]

	if strings.Contains(last, ".") {
		return strings.Split(last, ".")[0]
	}

	return last
}

// Path returns path of file
func (fp File) Path() string {
	if fp == "" {
		return ""
	}

	file := strings.Split(string(fp), "/")
	if len(file) == 1 {
		return file[0]
	}

	return strings.Join(file[:len(file)-2], "/")
}

func (fp File) SaveImagesByPattern(pattern string, images ...image.Image) error {
	es := ErrorStack{}

	for i, img := range images {
		val := strings.NewReplacer(
			"{{.Name}}", fp.Filename(),
			"{{.Counter}}", fmt.Sprint(i),
		).Replace(pattern)

		f, err := os.Create(val)
		if err != nil {
			es.Add(err)
			continue
		}
		defer f.Close()

		err = jpeg.Encode(f, img, &jpeg.Options{Quality: 75})
		es.Add(err)
	}

	return es.Get()
}
