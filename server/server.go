package server

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"html/template"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/StevenZack/gte/config"
	"github.com/StevenZack/gte/util"
)

type Server struct {
	HTTPServer           *http.Server
	cfg                  config.Config
	prehandlers          []func(w http.ResponseWriter, r *http.Request) bool
	funcs                template.FuncMap
	isRunningMode        bool //is in production mode
	precompiledTemplates *template.Template
}

func NewServer(cfg config.Config, isRunningMode bool) (*Server, error) {
	s := &Server{
		cfg:           cfg,
		isRunningMode: isRunningMode,
	}
	//funcs
	s.funcs = template.FuncMap{
		"httpGet":      s.httpGet,
		"httpGetJson":  s.httpGetJson,
		"mapOf":        util.MapOf,
		"httpPostJson": s.httpPostJson,
		"unescape":     unescape,
		"startsWith":   strings.HasPrefix,
		"endsWith":     strings.HasSuffix,
	}
	//route duplication check
	routeMap := map[string]string{}
	for _, route := range s.cfg.Routes {
		f := util.FormatParam(route.Path)
		exists, ok := routeMap[f]
		if ok {
			return nil, errors.New("Duplicate route path: '" + route.Path + "' with '" + exists + "'")
		}
		routeMap[f] = route.Path
	}

	s.HTTPServer = &http.Server{Addr: cfg.Host + ":" + strconv.Itoa(cfg.Port), Handler: s}

	// precompile in production mode
	if isRunningMode {
		var e error
		s.precompiledTemplates, e = util.ParseTemplates(s.cfg.Root, s.funcs)
		if e != nil {
			log.Println(e)
			return nil, e
		}
	}
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions && r.URL.Path == "/reload" {
		e := s.Reload()
		if e != nil {
			log.Println(e)
			http.Error(w, e.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("OK"))
		return
	}
	//prehandler
	for _, pre := range s.prehandlers {
		interrupt := pre(w, r)
		if interrupt {
			return
		}
	}

	//blacklist
	for _, black := range append(s.cfg.BlackList, s.cfg.InternalBlackList...) {
		if r.URL.Path == black {
			s.NotFound(w, r)
			return
		}
	}
	route := config.Route{
		Path: r.URL.Path,
		To:   r.URL.Path,
	}
	s.serveRoute(route, w, r, 0)
}

func (s *Server) serveRoute(route config.Route, w http.ResponseWriter, r *http.Request, statusCode int) {
	//route
	if route.To == "/" {
		route.To = "/index.html"
	}

	for _, cfgRoute := range s.cfg.Routes {
		if util.MatchRoute(cfgRoute.Path, r.URL.Path) {
			route.Path = cfgRoute.Path
			route.To = cfgRoute.To
		}
	}

	//lang
	ext := filepath.Ext(route.To)
	prefix := strings.TrimSuffix(route.To, ext)
	if _, e := os.Stat(filepath.Join(s.cfg.Root, prefix+"_"+util.GetLangShort(r)+ext)); e == nil {
		route.To = prefix + "_" + util.GetLangShort(r) + ext
	} else if _, e := os.Stat(filepath.Join(s.cfg.Root, prefix+"_"+util.GetLang(r)+ext)); e == nil {
		route.To = prefix + "_" + util.GetLang(r) + ext
	}

	//serve file
	switch ext {
	case ".html":
		w.Header().Set("Content-Type", "text/html")

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
		}

		if st, e := os.Stat(filepath.Join(s.cfg.Root, route.To)); e == nil {
			w.Header().Set("Last-Modified", st.ModTime().Format(http.TimeFormat))
		}
	default:
		path := filepath.Join(s.cfg.Root, route.To)
		if util.ShouldCWebp(ext) && strings.Contains(r.Header.Get("Accept"), "webp") {
			if _, e := os.Stat(path + ".webp"); e == nil {
				http.ServeFile(w, r, path+".webp")
				return
			}
		}

		//gzip
		if util.ShouldGZip(ext) && strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			if _, e := os.Stat(path + ".gzip"); e == nil {
				w.Header().Set("Content-Encoding", "gzip")
				w.Header().Set("Content-Type", mime.TypeByExtension(ext))
				http.ServeFile(w, r, path+".gzip")
				return
			}
		}

		if _, e := os.Stat(path); os.IsNotExist(e) {
			s.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, path)
		return
	}

	//parse templates
	var e error
	var t *template.Template

	if s.isRunningMode {
		t = s.precompiledTemplates
	} else {
		t, e = util.ParseTemplates(s.cfg.Root, s.funcs)
		if e != nil {
			log.Println(e)
			http.Error(w, e.Error(), http.StatusInternalServerError)
			return
		}
		if t == nil {
			log.Println("t == nil")
			if statusCode == 404 {
				http.NotFound(w, r)
				return
			}
			s.NotFound(w, r)
			return
		}
	}

	out := new(bytes.Buffer)
	e = t.ExecuteTemplate(out, route.To, NewContext(s.cfg, route, w, r))
	if e != nil {
		if strings.Contains(e.Error(), "is undefined") {
			log.Println(e)
			s.NotFound(w, r)
			return
		}

		log.Println(e)
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	//gzip

	if statusCode > 0 {
		w.WriteHeader(statusCode)
	}
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		rw := gzip.NewWriter(w)
		defer rw.Close()
		rw.Name, e = url.PathUnescape(filepath.Base(route.To))
		if e != nil {
			log.Println(e)
			http.Error(w, e.Error(), http.StatusInternalServerError)
			return
		}
		_, e = io.Copy(rw, out)
		if e != nil {
			log.Println(e)
			http.Error(w, e.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	w.Write(out.Bytes())
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
	if s.cfg.NotFoundPage != "" {
		s.serveRoute(config.Route{
			Path: r.URL.Path,
			To:   s.cfg.NotFoundPage,
		}, w, r, 404)
		return
	}
	http.NotFound(w, r)
}

func (s *Server) Reload() error {
	cfg, e := config.LoadConfig(s.cfg.Env, s.cfg.Root, s.cfg.Port)
	if e != nil {
		log.Println(e)
		return e
	}

	s.cfg = cfg

	if s.isRunningMode {
		s.precompiledTemplates, e = util.ParseTemplates(s.cfg.Root, s.funcs)
		if e != nil {
			log.Println(e)
			return e
		}

	}
	log.Println("server reloaded")
	return nil
}
