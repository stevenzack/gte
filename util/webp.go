package util

import (
	"os/exec"
)

func CWebp(file, out string) error {
	// return ("cwebp", "-o", out, file)
	return exec.Command("cwebp", "-o", out, file).Run()
}
