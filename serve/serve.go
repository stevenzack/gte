package serve

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/StevenZack/gte/util"
	"github.com/StevenZack/mux"
	"github.com/StevenZack/openurl"
	"github.com/urfave/cli"
)

var (
	server *mux.Server
	dir    string
)

func ApiCommand(c *cli.Context) error {
	return serve(c.String("dir"), c.Int("p"))
}

func serve(d string, port int) error {
	server = mux.NewServer("localhost:" + strconv.Itoa(port))
	dir = d

	//validate
	info, e := os.Stat(dir)
	if e != nil {
		log.Println(e)
		return e
	}
	if !info.IsDir() {
		return errors.New("'" + dir + "' is not a dir")
	}

	server.HandleMultiReqs("/", handler)

	fmt.Println("Running server on " + server.HTTPServer.Addr)
	openurl.Open("http://" + server.HTTPServer.Addr)
	e = server.ListenAndServe()
	if e != nil {
		if strings.Contains(e.Error(), "bind:") {
			return serve(d, port+1)
		}
		log.Println(e)
		return e
	}
	return nil
}
func handler(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	if uri == "/" {
		uri = "/index.html"
	}

	fmt.Print("\n" + uri)
	//serve file
	switch filepath.Ext(uri) {
	case ".html":
		fmt.Print("\t[T]")
	default:
		http.ServeFile(w, r, filepath.Join(dir, uri))
		return
	}

	//parse template
	t, e := util.ParseTemplates(dir)
	if e != nil {
		log.Println(e)
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	if t == nil {
		server.NotFound(w, r)
		return
	}

	e = t.ExecuteTemplate(w, uri, nil)
	if e != nil {
		log.Println(e)
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
}
