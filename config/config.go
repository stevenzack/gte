package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/StevenZack/gte/util"
	"github.com/StevenZack/tools/strToolkit"
)

type Config struct {
	Host              string            `json:"host"`
	Port              int               `json:"port"`
	Root              string            `json:"root"` //root directory of your project
	Routes            []Route           `json:"routes"`
	BlackList         []string          `json:"blackList"`
	ApiServer         string            `json:"apiServer"`
	Envs              map[string]Config `json:"envs"` //Environments
	Env               string            `json:"-"`
	InternalBlackList []string          `json:"-"`
}
type Route struct {
	Path string `json:"path"`
	To   string `json:"to"`
}

const (
	CONFIG_FILE_NAME = "gte.config.json"
)

func LoadConfig(env, root string, port int) (Config, error) {
	v := Config{
		Env:  env,
		Host: "0.0.0.0",
		Port: port,
		Root: root,
		InternalBlackList: []string{
			"/" + CONFIG_FILE_NAME,
		},
		ApiServer: "http://localhost",
	}

	b, e := ioutil.ReadFile(filepath.Join(root, CONFIG_FILE_NAME))
	if e != nil {
		if os.IsNotExist(e) {
			return v, nil
		}
		return v, e
	}

	e = json.Unmarshal(b, &v)
	if e != nil {
		return v, e
	}

	//handle envs
	if v.Envs != nil && env != "" {
		v1, ok := v.Envs[env]
		if !ok {
			return v, errors.New("No environment named '" + env + "'")
		}
		e := util.ReplaceFieldIND(&v, v1)
		if e != nil {
			return v, e
		}
	}
	return v, nil
}

func (r *Route) Params(uri string) map[string]string {
	ss1 := strings.Split(r.Path, "/")
	ss2 := strings.Split(uri, "/")
	m := make(map[string]string)
	for i := 0; i < len(ss1) && i < len(ss2); i++ {
		k := ss1[i]
		if k == "" {
			continue
		}
		if !strings.HasPrefix(k, ":") {
			continue
		}
		k = strToolkit.TrimStart(k, ":")

		v := ss2[i]
		m[k] = v
	}
	return m
}

func (r *Route) ParamPrefix() (string, bool) {
	i := strings.Index(r.Path, "/:")
	if i == -1 {
		return "", false
	}
	return r.Path[:i], true
}
