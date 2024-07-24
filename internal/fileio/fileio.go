package fileio

import (
	"gomaker/internal/logger"
	"os"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func WriteFile(fileName string, contents []byte) {
    err := os.WriteFile(fileName, contents, 0644)
    if err != nil {
		logger.Error("Write-File", fileName, "Failed to write file")
	}
}

