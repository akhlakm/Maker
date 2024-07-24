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
    Env 		string
    Install 	[]Func
    Run 		[]Func
    Setup 		[]Func
    Add 		[]Func
    Build 		[]Func
    Deploy 		[]Func
}


func LoadYAML(filename string) {
    logger.Trace("Load-YAML", filename)
    if fileio.FileExists(filename) {
        contents, err := os.ReadFile(filename)
        if err != nil {
            log.Fatal(err)
        }

        if err := yaml.Unmarshal(contents, &Defn); err != nil {
            logger.Error("Load-YAML", filename, "Failed to read YAML file.")
        }

        logger.Done("Load-YAML", filename)
        // fmt.Printf("%+v\n", Defn)
    }
}


func ListBlocks() {
    logger.Print("Available commands:")
    if len(Defn.Install) > 0 {
        logger.Print("  install   -   Install the application.")
    }
    if len(Defn.Run) > 0 {
        logger.Print("  run       -   Run the application.")
    }
    if len(Defn.Setup) > 0 {
        logger.Print("  setup     -   Prepare the development environment.")
    }
    if len(Defn.Add) > 0 {
        logger.Print("  add       -   Add new items to application.")
    }
    if len(Defn.Build) > 0 {
        logger.Print("  build     -   Build and test the application.")
    }
    if len(Defn.Deploy) > 0 {
        logger.Print("  deploy    -   Push or deploy the application.")
    }
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
    } else if block == "deploy" {
        printFunctions(block, Defn.Deploy)
    }
}

func buildFunctions(block string, flist []Func) string {
    // contents := fmt.Sprintf("\n## %s items #############", block)
    contents := ""
    for _, fitem := range flist {
        for name, body := range fitem {
            funcname := fmt.Sprintf("%s_%s", block, name)
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
    contents += buildFunctions("deploy", Defn.Deploy)

    contents += "\n\n" + entry + "\n"
    fileio.WriteFile(filename, []byte(contents))
}

func RunFunction(scriptpath string, block string, funcname string) {
    entry := fmt.Sprintf("%s_%s", block, funcname)
    saveMakeScript(scriptpath, entry)

    args := []string{scriptpath}
    args = append(args, os.Args[3:]...)

    // run from cwd
    ok, err := runner.Execute("/bin/bash", args, "")
    if !ok {
        logger.Error("Run-Function", entry, fmt.Sprintf("Failed: %s", err))
    }
}
