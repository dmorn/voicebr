package main

import (
	"flag"
	"log"
	"os"
	"net/http"
	"github.com/jecoz/voicebr"
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
	privateKey  = flag.String("key.pem", "private-key.pem", "Path to the private key that should be used to sign JWTs")
	appID  = flag.String("app.id", "test123", "Nexmo's application identifier")
	appNumber  = flag.String("app.number", "1111111111", "Nexmo's application registered number")
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

	file, err := os.Open(*privateKey)
	if err != nil {
		panic(err)
	}

	client, err := voicebr.NewClient(file, *appID, *appNumber, *hostAddr)
	file.Close()
	if err != nil {
		panic(err)
	}

	r := voicebr.NewRouter(client, *rootDir, *hostAddr)

	log.Printf("%v listening on port :%s", os.Args[0], *port)
	if err := http.ListenAndServe(":"+*port, r); err != nil {
		log.Fatal(err)
	}
}
