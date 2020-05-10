// Holds code for retrieving the default file paths.
// If we want to allow the client to change path,
// we can add the logic here.
package paths

import "path/filepath"

var (
	RootDir = rootDir
	Prefs   = prefs
)

func rootDir() string {
	return "."
}

func prefs() string {
	return filepath.Join(RootDir(), "voileyprefs.hujson")
}
