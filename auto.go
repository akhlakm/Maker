package main

import (
	"fmt"
	"gomaker/internal/fileio"
	"gomaker/internal/logger"
	"gomaker/internal/runner"
	"os"
	"path/filepath"
)

func Run(script string, args []string) {
	logger.Trace("Auto-Run", script)

	// Get the path of the executable.
	ex, err := os.Executable()
	if err != nil {
		logger.Error("Auto-Run", "", "Failed to get executable path")
		return
	}

	// Join the path of the script.
	scriptPath := filepath.Join(filepath.Dir(ex), script)

	// Check if the script exists.
	if !fileio.FileExists(scriptPath) {
		logger.Error("Auto-Run", scriptPath, "Failed to find script")
		return
	}

	// Run the script with the arguments.
	ok, err := runner.Execute(scriptPath, args, "")
	if !ok {
		logger.Error("Auto-Run", scriptPath, fmt.Sprintf("Failed: %s", err))
	}
}
