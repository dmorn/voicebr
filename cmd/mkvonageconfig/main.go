package main

import "os"
import "github.com/jecoz/callrelay/vonage"

func main() {
	vonage.WriteConfig(os.Stdout, new(vonage.Config))
}
