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
	port = flag.String("port", "4001", "Server listening port")
)

func main() {
	flag.Parse()

	log.Printf("version: %s, commit: %s, built at: %s\n\n", version, commit, buildTime)

	r := voicebr.NewRouter()

	log.Printf("%v listening on port :%s", os.Args[0], *port)
	if err := http.ListenAndServe(":"+*port, r); err != nil {
		log.Fatal(err)
	}
}
