package main

import (
	"flag"
	"github.com/jecoz/voicebr"
	"log"
	"net/http"
	"os"
)

// Version and BuildTime are filled in during build by the Makefile
var (
	version   = "N/A"
	commit    = "N/A"
	buildTime = "N/A"
)

var (
	port     = flag.String("port", "4001", "Server listening port")
	hostAddr = flag.String("host.addr", "http://d1f61c3e.ngrok.io", "Canonical address of the publicly available web server")
	rootDir  = flag.String("root.dir", "", "Storage root directory. Defaults to the current dir")
)

func main() {
	flag.Parse()

	log.Printf("version: %s, commit: %s, built at: %s\n\n", version, commit, buildTime)

	if *rootDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		*rootDir = wd
	}

	r := voicebr.NewRouter(*rootDir, *hostAddr)

	log.Printf("%v listening on port :%s", os.Args[0], *port)
	if err := http.ListenAndServe(":"+*port, r); err != nil {
		log.Fatal(err)
	}
}
