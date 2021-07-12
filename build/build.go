package build

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/StevenZack/gte/config"
	"github.com/StevenZack/gte/util"
	"github.com/StevenZack/tools/strToolkit"
	"github.com/urfave/cli"
)

const (
	DEFAULT_DESTINATION = "dist"
)

func ApiCommand(c *cli.Context) error {
	output := c.String("o")
	os.RemoveAll(output)
	e := os.MkdirAll(output, 0755)
	if e != nil {
		log.Println(e)
		return e
	}

	return build(c.Args().First(), output)
}

func build(env, out string) error {
	//validate
	info, e := os.Stat(out)
	if e != nil {
		log.Println(e)
		return e
	}
	if !info.IsDir() {
		return errors.New("'" + out + "' is not a directory")
	}

	root, e := filepath.Abs(".")
	if e != nil {
		log.Println(e)
		return e
	}
	root = filepath.ToSlash(root)

	e = filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		path = filepath.ToSlash(path)
		relativePath := strToolkit.TrimStart(path, root)
		if strings.Contains(relativePath, DEFAULT_DESTINATION+"/") || strings.Contains(relativePath, "/.") {
			return nil
		}
		dst := filepath.Join(out, relativePath)

		e := os.MkdirAll(filepath.Dir(dst), 0744)
		if e != nil {
			log.Println(e)
			return e
		}

		ext := filepath.Ext(info.Name())
		switch ext {
		case ".js":
			e := util.Obfuscate(path, dst)
			if e != nil {
				log.Println(e)
				return e
			}

		case ".css":
			e := util.MinifyCss(path, dst)
			if e != nil {
				log.Println(e)
				return e
			}
		default:
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
		}

		var tail string
		//gzip
		if util.ShouldGZip(ext) && info.Name() != config.CONFIG_FILE_NAME {
			e = util.Gzip(path, dst+".gzip")
			if e != nil {
				log.Println(e)
				return e
			}
			tail = "\t[gzip]"
		}

		if util.ShouldCWebp(ext) {
			e = util.CWebp(path, dst+".webp")
			if e != nil {
				log.Println(e)
				return e
			}
			tail = "\t[webp]"
		}

		fmt.Println(relativePath, tail)
		return nil
	})
	if e != nil {
		log.Println(e)
		return e
	}

	fmt.Println("\n build finished: ", out)
	return nil
}
