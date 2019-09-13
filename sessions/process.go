// +build !js !wasm

package sessions

import (
	"os"

	"github.com/bugsnag/panicwrap"
)

// Checks to see if this is the application process, as opposed to the process
// that monitors for panics
func isApplicationProcess() bool {
	// Application process is run first, and this will only have been set when
	// the monitoring process runs
	return "" == os.Getenv(panicwrap.DEFAULT_COOKIE_KEY)
}
