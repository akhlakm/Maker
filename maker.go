package main

import (
	"gomaker/internal/fileio"
	"gomaker/internal/logger"
	"gomaker/internal/parser"
	"os"
)

var makeFile string = "maker.yaml"
var makeScript string = "/tmp/maker.sh"

func main() {
	logger.Trace("main", "")

	if fileio.FileExists(makeFile) {
		parser.LoadYAML(makeFile)

		if len(os.Args) == 1 {
			parser.ListBlocks()

		} else if len(os.Args) == 2 {
			// List block functions
			parser.ListFunctions(os.Args[1])
		} else {
			// call the block function and pass args
			parser.RunFunction(makeScript, os.Args[1], os.Args[2])
		}
	} else {
		logger.Print("No makefile found. Please run 'init' to create one.")
	}

	parser.LoadYAML(makeFile)

	logger.Print(".")
}

