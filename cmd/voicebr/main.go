/// Broadcast voice messages to a set of recipients.
/// Copyright (C) 2019 Daniel Morandini (jecoz)
///
/// This program is free software: you can redistribute it and/or modify
/// it under the terms of the GNU General Public License as published by
/// the Free Software Foundation, either version 3 of the License, or
/// (at your option) any later version.
///
/// This program is distributed in the hope that it will be useful,
/// but WITHOUT ANY WARRANTY; without even the implied warranty of
/// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
/// GNU General Public License for more details.
///
/// You should have received a copy of the GNU General Public License
/// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
	port         = flag.String("port", "4001", "Server listening port")
	hostAddr     = flag.String("host.addr", "http://d1f61c3e.ngrok.io", "Canonical address of the publicly available web server")
	rootDir      = flag.String("root.dir", "", "Storage root directory. Defaults to the current dir")
	privateKey   = flag.String("key.pem", "private-key.pem", "Path to the private key that should be used to sign JWTs")
	appID        = flag.String("app.id", "test123", "Nexmo's application identifier")
	appNumber    = flag.String("app.number", "1111111111", "Nexmo's application registered number")
	contactsFile = flag.String("contacts.file", "data/contacts.csv", "Name of the input contacts to use when broadcasting calls.")
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

	s := &voicebr.Store{RootDir: *rootDir, ContactsFile: *contactsFile}
	r := voicebr.NewRouter(client, s, *hostAddr)

	log.Printf("%v listening on port :%s", os.Args[0], *port)
	if err := http.ListenAndServe(":"+*port, r); err != nil {
		log.Fatal(err)
	}
}
