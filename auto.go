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
	if os.Getenv("AUTODIR") != "" {
		return os.Getenv("AUTODIR")
	}

	// Get the path of the executable.
	ex, err := os.Executable()
	if err != nil {
		logger.Error("auto", "", "Failed to get executable path")
		return ""
	}

	// Return the directory of the executable as the script directory.
	return filepath.Dir(ex)
}


func main() {
	logger.Trace("auto", "")

	// Directory where all the scripts are located.
	script_dir := get_script_dir()

	if len(os.Args) < 2 {
		logger.Error("auto", script_dir, "No script provided")
		return
	}

	// script is the first argument.
	script := os.Args[1]

	// args are the rest of the arguments.
	args := os.Args[2:]

	
	// if current directory then we need to manually set the scriptpath
	// as filepath.Join will not add the prefix "./"
	var scriptPath string
	if script_dir == "." || script_dir == "./" {
		scriptPath = "./" + script
	} else {
		// Join the path of the script.
		scriptPath = filepath.Join(script_dir, script)
	}

	// Check if the script exists.
	if !fileio.FileExists(scriptPath) {
		logger.Error("auto", scriptPath, "Failed to find such script")
		return
	}

	// Run the script with the arguments.
	ok, err := runner.Execute(scriptPath, args, "")
	if !ok {
		logger.Error("auto", scriptPath, fmt.Sprintf("Failed: %s", err))
	}

	logger.Print("~")
}
