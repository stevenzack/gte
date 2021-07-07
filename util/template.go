package util

import (
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/StevenZack/tools/strToolkit"
)

func ParseTemplates(dir string, funcs template.FuncMap) (*template.Template, error) {
	abs, e := filepath.Abs(dir)
	if e != nil {
		return nil, e
	}
	abs = filepath.ToSlash(abs)

	var root *template.Template
	e = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		path = filepath.ToSlash(path)

		switch filepath.Ext(info.Name()) {
		case ".html":
			relativeUri := strToolkit.TrimStart(path, abs) // like /index.html
			if root == nil {
				root = template.New(relativeUri).Funcs(funcs)
			}

			var t *template.Template
			if relativeUri == root.Name() {
				t = root
			} else {
				t = root.New(relativeUri)
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

			_, e = t.Parse(string(b))
			if e != nil {
				return e
			}

		}
		return nil
	})
	if e != nil {
		return nil, e
	}

	return root, nil
}
