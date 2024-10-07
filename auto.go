package main

import (
	"fmt"
	"gomaker/internal/fileio"
	"gomaker/internal/logger"
	"gomaker/internal/runner"
	"os"
	"path/filepath"
)

func get_script_dir() string {
	// check if AUTODIR is set in the environment.
	var script_dir string
	if os.Getenv("AUTODIR") != "" {
		script_dir = os.Getenv("AUTODIR")
	} else {
		// Get the path of the executable.
		ex, err := os.Executable()
		if err != nil {
			logger.Error("auto", "", "Failed to get executable path")
			return ""
		}
		script_dir = filepath.Dir(ex)
	}

	// Get the absolute path of the script directory.
	abs_script_dir, err := filepath.Abs(script_dir)
	if err != nil {
		logger.Error("auto", "", "Failed to get absolute path of script directory")
	}

	// Set the AUTODIR environment variable.
	os.Setenv("AUTODIR", abs_script_dir)

	return abs_script_dir
}

func run_script(scriptPath string, args []string) {
	// Check if the script exists.
	if !fileio.FileExists(scriptPath) {
		logger.Error("auto", scriptPath, "Failed to find such script")
	}

	// Run the script with the arguments.
	ok, err := runner.Execute(scriptPath, args, "")
	if !ok {
		logger.Error("auto", scriptPath, fmt.Sprintf("Failed: %s", err))
	}
}

func main() {
	logger.Trace("auto", "")

	// Directory where all the scripts are located.
	script_dir := get_script_dir()

	if len(os.Args) >= 2 {

		// script is the first argument.
		script := os.Args[1]

		// args are the rest of the arguments.
		args := os.Args[2:]

		// run the script.
		scriptPath := filepath.Join(script_dir, script)
		run_script(scriptPath, args)

	} else {
		logger.Print("Usage: auto <script> [args...]\n")
		logger.Print("Scripts location: " + script_dir)
	}

	logger.Print("~")
}
