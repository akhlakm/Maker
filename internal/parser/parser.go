package parser

import (
	"fmt"
	"gomaker/internal/fileio"
	"gomaker/internal/logger"
	"gomaker/internal/runner"
	"log"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

type Func map[string]string

var Defn struct {
	Env     string
	Install []Func
	Run     []Func
	Setup   []Func
	Add     []Func
	Build   []Func
	Test    []Func
	Deploy  []Func
}

func LoadMakerYaml(filenames []string) bool {
	// Loop over the maker.yaml files and load them if any exists.
	// Returns true if a file was found, false otherwise.

	for _, filename := range filenames {
		if fileio.FileExists((filename)) {
			contents, err := os.ReadFile(filename)
			if err != nil {
				log.Fatal(err)
			}

			if err := yaml.Unmarshal(contents, &Defn); err != nil {
				logger.Error("Load-YAML", filename, "Failed to read YAML file.")
			}

			logger.Done("Load-YAML", filename)
			// fmt.Printf("%+v\n", Defn)

			return true
		}
	}

	return false
}

func functionNames(flist []Func) string {
	namelist := ""
	for _, fitem := range flist {
		for name := range fitem {
			if len(namelist) > 0 {
				namelist += ", " + name
			} else {
				namelist = name
			}
		}
	}
	return namelist
}

func ListBlocks() {
	logger.Print("Available block - commands:\n")
	if len(Defn.Install) > 0 {
		logger.Print(
			fmt.Sprintf("  install   -   %s", functionNames(Defn.Install)))
	}
	if len(Defn.Run) > 0 {
		logger.Print(
			fmt.Sprintf("  run       -   %s", functionNames(Defn.Run)))
	}
	if len(Defn.Setup) > 0 {
		logger.Print(
			fmt.Sprintf("  setup     -   %s", functionNames(Defn.Setup)))
	}
	if len(Defn.Add) > 0 {
		logger.Print(
			fmt.Sprintf("  add       -   %s", functionNames(Defn.Add)))
	}
	if len(Defn.Build) > 0 {
		logger.Print(
			fmt.Sprintf("  build     -   %s", functionNames(Defn.Build)))
	}
	if len(Defn.Test) > 0 {
		logger.Print(
			fmt.Sprintf("  test      -   %s", functionNames(Defn.Test)))
	}
	if len(Defn.Deploy) > 0 {
		logger.Print(
			fmt.Sprintf("  deploy    -   %s", functionNames(Defn.Deploy)))
	}

	logger.Print("\nPass a block name to see more details.")
}

func printFunctions(block string, flist []Func) {
	logger.Print(fmt.Sprintf("Available %s commands:", block))
	for _, fitem := range flist {
		for name, body := range fitem {
			summary := ""
			lines := strings.Split(body, "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "##") {
					summary += "\n        " + line
				}
			}
			logger.Print(fmt.Sprintf("\n  %s - %s", name, summary))
		}
	}

	logger.Print(fmt.Sprintf("\nPass %s <command> to execute it.", block))
}

func ListFunctions(block string) {
	if block == "env" {
		logger.Print(Defn.Env)
	} else if block == "install" {
		printFunctions(block, Defn.Install)
	} else if block == "run" {
		printFunctions(block, Defn.Run)
	} else if block == "setup" {
		printFunctions(block, Defn.Setup)
	} else if block == "add" {
		printFunctions(block, Defn.Add)
	} else if block == "build" {
		printFunctions(block, Defn.Build)
	} else if block == "test" {
		printFunctions(block, Defn.Test)
	} else if block == "deploy" {
		printFunctions(block, Defn.Deploy)
	} else {
		logger.Error("List-Function", block, "Unknown block.")
	}
}

func buildFunctions(block string, flist []Func) string {
	// contents := fmt.Sprintf("\n## %s items #############", block)
	contents := ""
	for _, fitem := range flist {
		for name, body := range fitem {
			funcname := fmt.Sprintf("%s-%s", block, name)
			function := fmt.Sprintf("function %s () {\n%s\nreturn 0\n}", funcname, body)
			contents += "\n\n" + function
		}
	}

	return contents
}

func saveMakeScript(filename string, entry string) {
	contents := "#!/bin/bash"
	contents += "\n\n" + Defn.Env + "\n"
	contents += buildFunctions("install", Defn.Install)
	contents += buildFunctions("run", Defn.Run)
	contents += buildFunctions("setup", Defn.Setup)
	contents += buildFunctions("add", Defn.Add)
	contents += buildFunctions("build", Defn.Build)
	contents += buildFunctions("test", Defn.Test)
	contents += buildFunctions("deploy", Defn.Deploy)

	contents += "\n\n" + entry + "\n"
	fileio.WriteFile(filename, []byte(contents))
}

func RunFunction(scriptpath string, block string, funcname string) {
	// call format: block-function "$@"
	entry := fmt.Sprintf("%s-%s \"$@\"", block, funcname)
	saveMakeScript(scriptpath, entry)

	// remove the calling arguments, pass the next ones to the script
	args := []string{scriptpath}
	args = append(args, os.Args[3:]...)

	// run from cwd
	ok, err := runner.Execute("/bin/bash", args, "")
	if !ok {
		logger.Error("Run-Function", entry, fmt.Sprintf("Failed: %s", err))
	}
}
