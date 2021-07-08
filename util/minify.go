package util

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/tdewolff/minify/v2/minify"
)

func MinifyCss(file, out string) error {
	b, e := ioutil.ReadFile(file)
	if e != nil {
		return e
	}
	result, e := minify.CSS(string(b))
	if e != nil {
		return e
	}
	e = ioutil.WriteFile(out, []byte(result), 0644)
	return e
}

func Gzip(file, out string) error {
	fo, e := os.OpenFile(out, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if e != nil {
		log.Println(e)
		return e
	}
	defer fo.Close()
	fi, e := os.OpenFile(file, os.O_RDONLY, 0644)
	if e != nil {
		log.Println(e)
		return e
	}
	defer fi.Close()
	writer := gzip.NewWriter(fo)
	writer.Name = filepath.Base(out)
	writer.Comment = "Gzip file of " + writer.Name
	writer.ModTime = time.Now()
	_, e = io.Copy(writer, fi)
	if e != nil {
		log.Println(e)
		return e
	}
	return nil
}
