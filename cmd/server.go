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

package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jecoz/voicebr/nexmo"
	"github.com/jecoz/voicebr/storage"
	"github.com/spf13/cobra"
)

var (
	appID   string
	appNum  string
	origin  string
	rootDir string
	pKey    string
	port    int
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start a voicebr server",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFlags(0)
		log.Printf("version: %s, commit: %s, built at: %s\n", Version, Commit, BuildTime)
		log.Printf("app-id: %s, app-num: %s, origin: %s, root-dir: %s\n\n", appID, appNum, origin, rootDir)

		log.Printf("loading private key from %s", pKey)
		file, err := os.Open(pKey)
		if err != nil {
			log.Fatal(err)
		}

		client, err := nexmo.NewClient(file, appID, appNum, origin)
		file.Close()
		if err != nil {
			panic(err)
		}

		log.Printf("creating local storage in: %s", rootDir)
		s := &storage.Local{RootDir: rootDir}
		r := nexmo.NewRouter(client, s, origin)

		log.Printf("%v listening on port :%d\n\n", os.Args[0], port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().IntVar(&port, "port", 4001, "Server listening port")
	serverCmd.Flags().StringVar(&rootDir, "root-dir", ".", "Root storage directory path")
	serverCmd.Flags().StringVar(&origin, "origin", "", "Canonical protocol + authority of the web server that will handle nexmo callbacks")
	serverCmd.Flags().StringVar(&pKey, "private-key", "", "Path to the private key that should be used to sign JWTs")
	serverCmd.Flags().StringVar(&appID, "app-id", "", "Nexmo's application identifier")
	serverCmd.Flags().StringVar(&appNum, "app-num", "", "Nexmo's application registered number")

	serverCmd.MarkFlagRequired("host-addr")
	serverCmd.MarkFlagRequired("app-id")
	serverCmd.MarkFlagRequired("app-num")
	serverCmd.MarkFlagRequired("private-key")
}
