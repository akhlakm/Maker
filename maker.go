package main

import (
	_ "embed"
	"fmt"
	"gomaker/internal/fileio"
	"gomaker/internal/logger"
	"gomaker/internal/parser"
	"os"
)

var makeFile string = "maker.yaml"
var makeScript string = "/tmp/maker.sh"

//go:embed sample.yaml
var sample []byte

func main() {
	logger.Trace("main", "")

	if fileio.FileExists(makeFile) {
		// makeFile found
		parser.LoadYAML(makeFile)

		if len(os.Args) == 1 {
			parser.ListBlocks()

		} else if len(os.Args) == 2 {
			// List functions of the block
			parser.ListFunctions(os.Args[1])
		} else {
			// call the block function
			parser.RunFunction(makeScript, os.Args[1], os.Args[2])
		}

	} else {
		// no makeFile

		if len(os.Args) == 2 && os.Args[1] == "init" {
			// write the sample as make file
			fileio.WriteFile(makeFile, sample)
		} else {
			logger.Print(
				fmt.Sprintf("No '%s' found. Please run 'init' to create one.",
					makeFile))
		}
	}

	parser.LoadYAML(makeFile)

	logger.Print(".")
}
