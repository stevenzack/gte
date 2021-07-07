package main

import (
	"log"
	"os"

	"github.com/StevenZack/gte/build"
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
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		{
			Name:  "serve",
			Usage: "Run a local server for developing",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "p",
					Usage: "specific server port",
				},
				cli.StringFlag{
					Name:  "dir",
					Usage: "Project root location",
				},
			},
			Action: serve.ApiCommand,
		},
		{
			Name:  "build",
			Usage: "Build HTML/CSS/JS file for production, including minifing, gzipping",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dir",
					Usage: "Project root location",
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
				},
				cli.StringFlag{
					Name:  "dir",
					Usage: "Project root location",
				},
			},
			Action: run.ApiCommand,
		},
	}

	e := app.Run(os.Args)
	if e != nil {
		log.Println(e)
		return
	}
}
