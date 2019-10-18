package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// CtxLogger serves as the base logger all other packages will use
var CtxLogger = log.New()

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	CtxLogger.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	CtxLogger.SetLevel(log.DebugLevel)
}
