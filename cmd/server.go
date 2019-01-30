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
	"log"
	"os"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/jecoz/voicebr/nexmo"
	"github.com/jecoz/voicebr/storage"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("version: %s, commit: %s, built at: %s\n\n", Version, Commit, BuildTime)

		file, err := os.Open(viper.GetString("private-key"))
		if err != nil {
			log.Fatal(err)
		}

		client, err := nexmo.NewClient(file, viper.GetString("app-id"), viper.GetString("app-num"), viper.GetString("host-addr"))
		file.Close()
		if err != nil {
			panic(err)
		}

		s := &storage.Local{RootDir: viper.GetString("root-dir")}
		r := nexmo.NewRouter(client, s, viper.GetString("host-addr"))

		log.Printf("%v listening on port :%d", os.Args[0], viper.GetInt("port"))
		if err := http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("port")), r); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().Int("port", 4001, "Server listening port")
	serverCmd.Flags().String("root-dir", ".", "Root storage directory path")
	serverCmd.Flags().String("host-addr", "", "Canonical address of the publicly available web server")
	serverCmd.Flags().String("private-key", "", "Path to the private key that should be used to sign JWTs")
	serverCmd.Flags().String("app-id", "", "Nexmo's application identifier")
	serverCmd.Flags().String("app-num", "", "Nexmo's application registered number")

	serverCmd.MarkFlagRequired("host-addr")
	serverCmd.MarkFlagRequired("private-key")
	serverCmd.MarkFlagRequired("app-id")
	serverCmd.MarkFlagRequired("app-num")

	viper.BindPFlag("port", serverCmd.Flags().Lookup("port"))
	viper.BindPFlag("host-addr", serverCmd.Flags().Lookup("host-addr"))
	viper.BindPFlag("root-dir", serverCmd.Flags().Lookup("root-dir"))
	viper.BindPFlag("private-key", serverCmd.Flags().Lookup("private-key"))
	viper.BindPFlag("app-id", serverCmd.Flags().Lookup("app-id"))
	viper.BindPFlag("app-num", serverCmd.Flags().Lookup("app-num"))
}
