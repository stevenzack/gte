package minimize

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli"
)

func ApiCommand(c *cli.Context) error {
	totalCount := 0
	e := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		switch filepath.Ext(path) {
		case ".mp4", ".png", ".jpeg", ".jpg":
			totalCount++
		case ".webp", ".gzip":
			fmt.Println(path, "!! deleted !!")
			os.Remove(path)
		default:
			return nil
		}
		return nil
	})
	if e != nil {
		log.Println(e)
		return e
	}

	count := 0
	var total, total2 int64
	e = filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		before := info.Size()
		var after int64
		switch filepath.Ext(path) {
		case ".mp4":
			fmt.Print(count, "/", totalCount, "\t", path)

			dst := path + ".min.mp4"
			e := exec.Command("ffmpeg", "-i", path, "-crf", "40", dst).Run()
			if e != nil {
				log.Println(e)
				return e
			}

			after, e = statSize(dst)
			if e != nil {
				log.Println(e)
				return e
			}

			e = exec.Command("mv", dst, path).Run()
			if e != nil {
				log.Println(e)
				return e
			}
		case ".png", ".jpeg", ".jpg":
			fmt.Print(count, "/", totalCount, "\t", path)

			dst := path + ".min.png"
			e := exec.Command("ffmpeg", "-i", path, dst).Run()
			if e != nil {
				log.Println(e)
				return e
			}
			after, e = statSize(dst)
			if e != nil {
				log.Println(e)
				return e
			}
			e = exec.Command("mv", dst, path).Run()
			if e != nil {
				log.Println(e)
				return e
			}
		default:
			return nil
		}

		fmt.Printf("\t[%d %%]\n", int((before-after)*100/before))
		count++
		total += before
		total2 += after
		return nil
	})
	if e != nil {
		log.Println(e)
		return e
	}

	fmt.Printf("compressed %d:\t compress rate %d %%\n", count, int((total-total2)*100/total))
	return nil
}

func statSize(path string) (int64, error) {
	info, e := os.Stat(path)
	if e != nil {
		return 0, e
	}
	return info.Size(), nil
}
