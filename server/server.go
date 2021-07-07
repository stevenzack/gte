package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/StevenZack/gte/config"
	"github.com/StevenZack/gte/util"
)

type Server struct {
	HTTPServer  *http.Server
	cfg         config.Config
	prehandlers []func(w http.ResponseWriter, r *http.Request) bool
}

func NewServer(cfg config.Config) (*Server, error) {
	s := &Server{
		cfg: cfg,
	}
	//route duplication check
	checked := map[string]string{}
	for _, route := range s.cfg.Routes {
		f := util.FormatParam(route.Path)
		exists, ok := checked[f]
		if ok {
			return nil, errors.New("Duplicate route path: '" + route.Path + "' with '" + exists + "'")
		}
		checked[f] = route.Path
	}

	s.HTTPServer = &http.Server{Addr: cfg.Host + ":" + strconv.Itoa(cfg.Port), Handler: s}
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//prehandler
	for _, pre := range s.prehandlers {
		interrupt := pre(w, r)
		if interrupt {
			return
		}
	}

	//blacklist
	for _, black := range s.cfg.BlackList {
		if r.URL.Path == black {
			s.NotFound(w, r)
			return
		}
	}

	//route
	route := config.Route{
		Path: r.URL.Path,
		To:   r.URL.Path,
	}
	if route.To == "/" {
		route.To = "/index.html"
	}

	for _, cfgRoute := range s.cfg.Routes {
		if util.MatchRoute(cfgRoute.Path, r.URL.Path) {
			route.Path = cfgRoute.Path
			route.To = cfgRoute.To
		}
	}

	//serve file
	switch filepath.Ext(route.To) {
	case ".html":
	default:
		http.ServeFile(w, r, filepath.Join(s.cfg.Root, route.To))
		return
	}

	//parse templates
	t, e := util.ParseTemplates(s.cfg.Root)
	if e != nil {
		log.Println(e)
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
	if t == nil {
		s.NotFound(w, r)
		return
	}

	e = t.ExecuteTemplate(w, route.To, NewContext(s.cfg, route, w, r))
	if e != nil {
		log.Println(e)
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) ListenAndServe() error {
	return s.HTTPServer.ListenAndServe()
}

func (s *Server) Stop() error {
	if s != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.
		e := s.HTTPServer.Shutdown(ctx)
		return e
	}
	return nil
}

func (s *Server) AddPrehandler(fn func(w http.ResponseWriter, r *http.Request) bool) {
	s.prehandlers = append(s.prehandlers, fn)
}

func (s *Server) NotFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}