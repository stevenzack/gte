package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type JsonResponse struct {
	StatusCode int
	Data       map[string]interface{}
	Error      string
}

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

func (s *Server) httpGetJson(url string) (*JsonResponse, error) {
	url = s.handleUrl(url)
	fmt.Println("GET\t", url)
	res, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()
	b, e := io.ReadAll(res.Body)
	if e != nil {
		return nil, e
	}

	rp := &JsonResponse{StatusCode: res.StatusCode}
	if res.StatusCode == http.StatusOK {
		v := make(map[string]interface{})
		e = json.Unmarshal(b, &v)
		if e != nil {
			return nil, e
		}
		rp.Data = v
	} else {
		rp.Error = string(b)
	}

	return rp, nil
}

func (s *Server) httpPostJson(url string, body interface{}) (*JsonResponse, error) {
	url = s.handleUrl(url)
	fmt.Println("POST\t", url)
	var reader *bytes.Reader
	if body != nil {
		b, e := json.Marshal(body)
		if e != nil {
			return nil, e
		}
		reader = bytes.NewReader(b)
	}

	res, e := http.Post(url, "application/json", reader)
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()

	b, e := io.ReadAll(res.Body)
	if e != nil {
		return nil, e
	}

	rp := &JsonResponse{StatusCode: res.StatusCode}
	if res.StatusCode == http.StatusOK {
		v := make(map[string]interface{})
		e = json.Unmarshal(b, &v)
		if e != nil {
			return nil, e
		}
		rp.Data = v
	} else {
		rp.Error = string(b)
	}

	return rp, nil
}

func unescape(s string) template.HTML {
	return template.HTML(s)
}
