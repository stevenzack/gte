package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/StevenZack/gte/config"
	"github.com/StevenZack/tools/strToolkit"
	"golang.org/x/text/language"
)

type Context struct {
	Config   config.Config
	route    config.Route
	Request  *Request
	Response *Response
}

func NewContext(cfg config.Config, route config.Route, w http.ResponseWriter, r *http.Request) *Context {
	ctx := &Context{
		Config: cfg,
		route:  route,
	}
	ctx.Request = NewRequest(ctx, r)
	ctx.Response = NewResponse(ctx, w)
	return ctx
}

func (c *Context) GetStr(key string) (string, error) {
	if c.Config.Lang.Dir == "" {
		return "", errors.New("Calling .GetStr() function, but 'lang' config is not set in  'gte.config.json' file")
	}
	accept := c.Request.Header.Get("Accept-Language")
	accept = strToolkit.SubBefore(accept, ";", accept)
	accept = strToolkit.SubBefore(accept, ",", accept)
	if accept == "" {
		return key, nil
	}

	tag, e := language.Parse(accept)
	if e != nil {
		return "", e
	}

	lang := tag.String()
	m, ok := c.Config.Strs[tag.String()]
	if !ok {
		base, _ := tag.Base()
		m, ok = c.Config.Strs[base.String()]
		if !ok {

			//return key as value when request of default language comes
			if (lang == c.Config.Lang.Default || base.String() == c.Config.Lang.Default) && c.Config.Lang.KeyAsValue {
				return key, nil
			}

			log.Println("Unsupported language '" + tag.String() + "', using default '" + c.Config.Lang.Default + "'")
			m, ok = c.Config.Strs[c.Config.Lang.Default]
			lang = c.Config.Lang.Default
		}
	}
	v, ok := m[key]
	if !ok {
		if c.Config.Lang.KeyAsValue {
			return key, nil
		}
		return "", errors.New("translation for key '" + key + "' not found in language resource file '" + lang + ".yaml'")
	}
	return v, nil

}
