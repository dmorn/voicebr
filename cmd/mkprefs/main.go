package main

import "os"
import "github.com/jecoz/callrelay"

func main() {
	callrelay.WritePrefs(os.Stdout, new(callrelay.Prefs))
}
