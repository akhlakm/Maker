package main

import (
	_ "embed"
	"fmt"
	"gomaker/internal/fileio"
	"gomaker/internal/logger"
	"gomaker/internal/parser"
	"os"
)

var makeScript string = "/tmp/maker.sh"

// List of possible maker file names
// Order of preference is maintained
var makeFiles = []string{"maker.yaml", "maker.yml"}

//go:embed sample.yml
var sample []byte

func main() {
	logger.Trace("main", "")

	if parser.LoadMakerYaml(makeFiles) {
		// makeFile found

		if len(os.Args) == 1 {
			parser.ListBlocks()

		} else if len(os.Args) == 2 {
			if os.Args[1] == "export" {
				parser.ExportMakerfile("maker.sh")
			} else {
				// List functions of the block
				parser.ListFunctions(os.Args[1])
			}
		} else {
			// call the block function
			parser.RunFunction(makeScript, os.Args[1], os.Args[2])
		}

	} else {
		// no makeFile

		// default file name for init
		makeFile := makeFiles[0]

		if len(os.Args) == 2 && os.Args[1] == "init" {
			// write the sample as make file
			fileio.WriteFile(makeFile, sample)
		} else {
			logger.Print(
				fmt.Sprintf("No '%s' found. Please run 'init' to create one.",
					makeFile))
		}
	}

	logger.Print("~")
}
