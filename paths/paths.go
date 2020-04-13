// Holds code for retrieving the default file paths.
//
// Internal API subject to change.
// If we want to allow the client to change
// path, we can add the logic here.
package paths

import "path/filepath"

var RootDir func() string

func rootDir() string {
	if f := RootDir; f != nil {
		return f()
	}
	return defaultRootDir()
}

// Path to user defined preferences.
func PrefsPath() string {
	return filepath.Join(rootDir(), "voicebr.prefs.hujson")
}

// Path to Vonage's configuration.
func VonageConfigPath() string {
	return filepath.Join(rootDir(), "vonage.config.hujson")
}

// Path to the file that holds the server's
// pid. The server will not start unless the file
// is deleted.
func ServerPidPath() string {
	return filepath.Join(rootDir(), "server.pid")
}
