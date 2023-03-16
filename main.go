package main

import (
	"log"
	"os"

	"github.com/StevenZack/gte/build"
	"github.com/StevenZack/gte/minimize"
	"github.com/StevenZack/gte/reload"
	"github.com/StevenZack/gte/run"
	"github.com/StevenZack/gte/serve"
	"github.com/urfave/cli"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	app := cli.NewApp()
	app.Name = "Golang Template Engine"
	app.Version = "1.2.0"
	wd, e := os.Getwd()
	if e != nil {
		log.Println(e)
		return
	}

	app.Commands = []cli.Command{
		{
			Name:  "serve",
			Usage: "Run a local server for developing",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "p",
					Usage: "specific server port",
					Value: 8080,
				},
				cli.StringFlag{
					Name:  "dir",
					Usage: "Project root location",
					Value: wd,
				},
				cli.StringFlag{
					Name:  "c",
					Usage: "GTE config json file location, default 'gte.config.json'",
					Value: "gte.config.json",
				},
			},
			Action: serve.ApiCommand,
		},
		{
			Name:  "build",
			Usage: "Build HTML/CSS/JS file for production, including minifing, gzipping",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "o",
					Usage: "Target build output location",
					Value: build.DEFAULT_DESTINATION,
				},
			},
			Action: build.ApiCommand,
		},
		{
			Name:  "run",
			Usage: "Run a production server",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "p",
					Usage: "specific server port",
					Value: 8080,
				},
				cli.StringFlag{
					Name:  "dir",
					Usage: "Project root location",
					Value: wd,
				},
				cli.StringFlag{
					Name:  "c",
					Usage: "GTE config json file location, default 'gte.config.json'",
					Value: "gte.config.json",
				},
			},
			Action: run.ApiCommand,
		},
		{
			Name:  "reload",
			Usage: "reload configure on server",
			Flags: []cli.Flag{
				cli.Int64Flag{
					Name:  "p",
					Usage: "port of server that you want to reload",
				},
			},
			Action: reload.ApiCommand,
		},
		{
			Name:   "minimize",
			Usage:  "minimize video,png,jpeg in current directory",
			Action: minimize.ApiCommand,
		},
	}

	e = app.Run(os.Args)
	if e != nil {
		log.Println(e)
		return
	}
}
