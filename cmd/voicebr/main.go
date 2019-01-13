package main

import (
	"log"
)

// Version and BuildTime are filled in during build by the Makefile
var (
	version   = "N/A"
	commit    = "N/A"
	buildTime = "N/A"
)

func main() {
	log.Printf("version: %s, commit: %s, built at: %s\n\n", version, commit, buildTime)
}
