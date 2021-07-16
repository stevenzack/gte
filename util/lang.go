package util

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/StevenZack/tools/strToolkit"
	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v2"
)

const (
	LANG_FILE_EXT = ".json"
)

func LoadJsonLangFile(path string) (map[string]string, error) {
	b, e := ioutil.ReadFile(path)
	if e != nil {
		log.Println(e)
		return nil, e
	}

	m := make(map[string]string)
	e = json.Unmarshal(b, &m)
	if e != nil {
		log.Println(e)
		return nil, e
	}

	out := make(map[string]string)
	needReplace := false
	for k, v := range m {
		k2 := upperCase(k)
		out[k2] = v
		if k2 != k {
			needReplace = true
		}
	}
	if needReplace {
		b, e = json.MarshalIndent(out, "", "\t")
		if e != nil {
			log.Println(e)
			return nil, e
		}
		e = ioutil.WriteFile(path, b, 0644)
		if e != nil {
			log.Println(e)
			return nil, e
		}
	}

	return out, nil
}

func FormatJsonLangFile(root, path string) (map[string]string, error) {
	b, e := ioutil.ReadFile(path)
	if e != nil {
		log.Println(e)
		return nil, e
	}

	m := make(map[string]string)
	e = yaml.Unmarshal(b, &m)
	if e != nil {
		log.Println(e)
		return nil, e
	}

	out := make(map[string]string)
	replacement := make(map[string]string)
	needReplace := false
	for k, v := range m {
		k2 := upperCase(k)
		out[k2] = v
		if k2 != k {
			needReplace = true
			replacement[k] = k2
		}
	}
	if needReplace {
		b, e = json.MarshalIndent(out, "", "\t")
		if e != nil {
			log.Println(e)
			return nil, e
		}
		e = ioutil.WriteFile(path, b, 0644)
		if e != nil {
			log.Println(e)
			return nil, e
		}

		e = filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
			if filepath.Ext(info.Name()) != ".html" {
				return nil
			}
			contentBytes, e := ioutil.ReadFile(path)
			if e != nil {
				log.Println(e)
				return e
			}
			content := string(contentBytes)

			for k, k2 := range replacement {
				content = strings.ReplaceAll(content, ".GetStr \""+k+"\"", ".GetStr \""+k2+"\"")
			}
			e = ioutil.WriteFile(path, []byte(content), 0644)
			if e != nil {
				log.Println(e)
				return e
			}
			return nil
		})
		if e != nil {
			log.Println(e)
			return nil, e
		}
	}

	return out, nil
}

func LoadYamlLangFile(path string) ([]byte, error) {
	b, e := ioutil.ReadFile(path)
	if e != nil {
		log.Println(e)
		return nil, e
	}
	content := string(b)
	if strings.Contains(content, "：") {
		content = strings.ReplaceAll(content, "：", ": ")
		content = strings.ReplaceAll(content, "＃", "#")
	}

	reader := bufio.NewReader(strings.NewReader(content))

	out := new(bytes.Buffer)
	for {
		s, e := reader.ReadString('\n')
		if e != nil {
			if e == io.EOF {
				break
			}
			log.Println(e)
			return nil, e
		}
		line := strings.TrimSuffix(s, "\n")
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			out.WriteString(line)
		} else {
			key := strToolkit.SubBefore(line, ": ", line)
			key = upperCase(unquote(key))
			value := strToolkit.SubAfter(line, ": ", "")
			value = quote(value)

			out.WriteString(key + ": " + value)
		}

		out.WriteString("\n")
	}

	e = ioutil.WriteFile(path, out.Bytes(), 0644)
	if e != nil {
		log.Println(e)
		return nil, e
	}

	return out.Bytes(), nil
}
func unquote(s string) string {
	s = strings.TrimPrefix(s, "\"")
	s = strings.TrimSuffix(s, "\"")
	return s
}
func quote(s string) string {
	if !strings.HasPrefix(s, "\"") {
		s = "\"" + s
	}
	if !strings.HasSuffix(s, "\"") {
		s = s + "\""
	}
	return s
}

func upperCase(s string) string {
	s = strings.ToUpper(strcase.ToSnake(s))
	s = strings.TrimSuffix(s, "_")
	return s + "_"
}
