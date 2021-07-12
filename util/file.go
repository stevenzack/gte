package util

import (
	"io/ioutil"
	"log"
)

func CopyFile(path, dst string) error {
	b, e := ioutil.ReadFile(path)
	if e != nil {
		log.Println(e)
		return e
	}
	e = ioutil.WriteFile(dst, b, 0644)
	if e != nil {
		log.Println(e)
		return e
	}
	return nil
}
