// Holds code for retrieving the default file paths.
//
// Internal API subject to change.
// If we want to allow the client to change
// path, we can add the logic here.
package paths

// Path to user preferences, such as NCCOs.
func PrefsPath() string {
	return "var/lib/voicebr/prefs.hujson"
}

// Path to Vonage's configuration.
func VonageConfigPath() string {
	return "var/lib/voicebr/vonage.hujson"
}
