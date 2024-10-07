package main

import (
	"gomaker/internal/logger"
)

func main() {
	logger.Trace("main", "")
	logger.Print("Hello, World!")
	logger.Print("~")
}
