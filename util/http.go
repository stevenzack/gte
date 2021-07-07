package util

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

func httpGet(url string) (string, error) {
	res, e := http.Get(url)
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

func httpGetJson(url string) (map[string]interface{}, error) {
	res, e := http.Get(url)
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
