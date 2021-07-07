package util

import (
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/StevenZack/tools/strToolkit"
)

func ParseTemplates(dir string) (*template.Template, error) {
	abs, e := filepath.Abs(dir)
	if e != nil {
		return nil, e
	}

	var t *template.Template
	e = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		switch filepath.Ext(info.Name()) {
		case ".html":
			relativeUri := strToolkit.TrimStart(path, abs) // like /index.html
			if t == nil {
				t = template.New(relativeUri)
			} else {
				t = t.New(relativeUri)
			}

			//read
			fi, e := os.OpenFile(path, os.O_RDONLY, 0644)
			if e != nil {
				return e
			}
			defer fi.Close()
			b, e := io.ReadAll(fi)
			if e != nil {
				return e
			}

			t, e = t.Parse(string(b))
			if e != nil {
				return e
			}

		}
		return nil
	})
	if e != nil {
		return nil, e
	}

	return t, nil
}
