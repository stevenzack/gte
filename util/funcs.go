package util

import "html/template"

var (
	Funcs = template.FuncMap{
		"httpGet":     httpGet,
		"httpGetJson": httpGetJson,
	}
)
