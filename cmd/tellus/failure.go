package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

var (
	recoveryFile    = "tellus.panic.log"
	recoveryMessage = `################################# TELLUS CRASH #################################

Tellus has recovered from a panic but had to shut down ungracefully. This is
always a bug! If you would, please file the output placed in the local directory
at "` + recoveryFile + `" as an issue at github.com/asteris-llc/tellus, along
with what you were trying to do at the time of the crash. Thanks, and sorry for
the hassle!

################################# TELLUS CRASH #################################`
)

func Recovery() {
	if r := recover(); r != nil {
		message := fmt.Sprintf("%s\n---\n%s", Version, r)
		// first print the recovery message and status to the screen, just in case
		// dumping fails.
		fmt.Println(message)
		fmt.Println("\n" + recoveryMessage)

		// now write to the disk! Ignoring the error intentionally because we've
		// already written to the screen.
		_ = ioutil.WriteFile(recoveryFile, []byte(message), 0644)
		os.Exit(2)
	}
}

// GracefullyFail handles cleanup for Tellus. Right now it doesn't have to do
// much, it just prints the given error message to the screen and exits. This
// function should handle all cleanup necessary, however.
func GracefullyFail(message string) {
	fmt.Println(message)
	os.Exit(1)
}
