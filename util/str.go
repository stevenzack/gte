package util

import (
	"net/http"
	"strings"

	"github.com/StevenZack/tools/strToolkit"
)

func Contains(s, sep string) bool {
	return strings.Contains(s, sep)
}

func GetLang(r *http.Request) string {
	accept := r.Header.Get("Accept-Language")
	accept = strToolkit.SubBefore(accept, ";", accept)
	accept = strToolkit.SubBefore(accept, ",", accept)
	return accept
}

func GetLangShort(r *http.Request) string {
	lang := GetLang(r)
	lang = strToolkit.SubBefore(lang, "-", lang)
	lang = strToolkit.SubBefore(lang, "_", lang)
	return lang
}
