package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
)

type JsonResponse struct {
	StatusCode int
	Data       map[string]interface{}
	Error      string
}

type StringResponse struct {
	StatusCode int
	Data       string
	Error      string
}

func (s *Server) handleUrl(url string) string {
	if strings.HasPrefix(url, "http") {
		return url
	}
	return s.cfg.ApiServer + url
}

func (s *Server) httpGet(url string) (*StringResponse, error) {
	res, e := http.Get(s.handleUrl(url))
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()
	b, e := io.ReadAll(res.Body)
	if e != nil {
		return nil, e
	}

	rp := StringResponse{
		StatusCode: res.StatusCode,
	}
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		rp.Data = string(b)
	} else {
		rp.Error = string(b)
	}

	return &rp, nil
}

func (s *Server) httpGetJson(url string) (*JsonResponse, error) {
	url = s.handleUrl(url)
	res, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()
	b, e := io.ReadAll(res.Body)
	if e != nil {
		return nil, e
	}
	fmt.Println("GET\t", res.StatusCode, "\t", url)

	rp := &JsonResponse{StatusCode: res.StatusCode}
	if res.StatusCode == http.StatusOK {
		v := make(map[string]interface{})
		e = json.Unmarshal(b, &v)
		if e != nil {
			rp.Error = string(e.Error())
		} else {
			rp.Data = v
		}
	} else {
		rp.Error = string(b)
		if strings.HasPrefix(rp.Error, "{") {
			v := make(map[string]interface{})
			if e = json.Unmarshal(b, &v); e == nil {
				rp.Data = v
			}
		}
	}

	return rp, nil
}

func (s *Server) httpPostJson(url string, body interface{}) (*JsonResponse, error) {
	url = s.handleUrl(url)
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

	fmt.Println("POST\t", res.StatusCode, "\t", url)
	b, e := io.ReadAll(res.Body)
	if e != nil {
		return nil, e
	}

	rp := &JsonResponse{StatusCode: res.StatusCode}
	if res.StatusCode == http.StatusOK {
		v := make(map[string]interface{})
		e = json.Unmarshal(b, &v)
		if e != nil {
			rp.Error = string(e.Error())
		} else {
			rp.Data = v
		}
	} else {
		rp.Error = string(b)
		if strings.HasPrefix(rp.Error, "{") {
			v := make(map[string]interface{})
			if e = json.Unmarshal(b, &v); e == nil {
				rp.Data = v
			}
		}
	}

	return rp, nil
}

func unescape(s string) template.HTML {
	return template.HTML(s)
}
