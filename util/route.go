package util

import "strings"

func FormatParam(path string) string {
	ss := strings.Split(path, "/")
	out := []string{}
	for _, s := range ss {
		if strings.HasPrefix(s, ":") {
			s = "-"
		}
		out = append(out, s)
	}
	return strings.Join(out, "/")
}

func MatchInParam(s1, s2 string) bool {
	return FormatParam(s1) == FormatParam(s2)
}

func MatchRoute(path, url string) bool {
	ss1 := strings.Split(path, "/")
	ss2 := strings.Split(url, "/")
	if len(ss1) != len(ss2) {
		return false
	}

	for i, s := range ss1 {
		if i == 0 {
			continue
		}
		v := ss2[i]
		if strings.HasPrefix(s, ":") {
			if v == "" {
				return false
			}
			continue
		}
		if v == "" {
			return false
		}
		if s != v {
			return false
		}
	}
	return true
}
