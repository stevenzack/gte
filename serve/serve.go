package serve

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/StevenZack/gte/config"
	"github.com/StevenZack/gte/server"
	"github.com/StevenZack/openurl"
	"github.com/urfave/cli"
)

func ApiCommand(c *cli.Context) error {
	return serve(c.Args().First(), c.String("dir"), c.Int("p"))
}

func serve(env, dir string, port int) error {
	//validate
	info, e := os.Stat(dir)
	if e != nil {
		log.Println(e)
		return e
	}
	if !info.IsDir() {
		return errors.New("'" + dir + "' is not a dir")
	}

	cfg, e := config.LoadConfig(env, dir, port)
	if e != nil {
		log.Println(e)
		return e
	}

	server, e := server.NewServer(cfg)
	if e != nil {
		log.Println(e)
		return e
	}
	server.AddPrehandler(printRequest)

	fmt.Println("Running server on " + server.HTTPServer.Addr)
	openurl.Open("http://" + server.HTTPServer.Addr)
	e = server.ListenAndServe()
	if e != nil {
		if strings.Contains(e.Error(), "bind:") {
			return serve(env, dir, port+1)
		}
		log.Println(e)
		return e
	}
	return nil
}

func printRequest(w http.ResponseWriter, r *http.Request) bool {
	fmt.Println(r.URL.Path)
	return false
}
