package main

import "os"
import "github.com/jecoz/voicebr"

func main() {
	voicebr.WritePrefs(os.Stdout, new(voicebr.Prefs))
}
