package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) handleUrl(url string) string {
	if strings.HasPrefix(url, "http") {
		return url
	}
	return s.cfg.ApiServer + url
}

func (s *Server) httpGet(url string) (string, error) {
	res, e := http.Get(s.handleUrl(url))
	if e != nil {
		return "", e
	}
	defer res.Body.Close()
	b, e := io.ReadAll(res.Body)
	if e != nil {
		return "", e
	}

	if res.StatusCode != http.StatusOK {
		return "", errors.New(strconv.Itoa(res.StatusCode) + ":" + string(b))
	}

	return string(b), nil
}

func (s *Server) httpGetJson(url string) (map[string]interface{}, error) {
	res, e := http.Get(s.handleUrl(url))
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()
	b, e := io.ReadAll(res.Body)
	if e != nil {
		return nil, e
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(strconv.Itoa(res.StatusCode) + ":" + string(b))
	}

	v := make(map[string]interface{})
	e = json.Unmarshal(b, &v)
	if e != nil {
		return nil, e
	}

	return v, nil
}
