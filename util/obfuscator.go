package util

import "github.com/StevenZack/tools/ioToolkit"

func Obfuscate(file, out string) error {
	return ioToolkit.RunAttachedCmd("javascript-obfuscator", file,
		"--output", out,
		"--control-flow-flattening", "true",
		"--control-flow-flattening-threshold", "1",
		"--numbers-to-expressions", "true",
		"--simplify", "true",
		"--shuffle-string-array", "true",
		"--split-strings", "true",
		"--string-array-threshold", "1",
	)
}
