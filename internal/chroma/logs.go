package chroma

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// Setup log formatting
func (chroma *Chroma) setupLogging() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "02-01-2006 15:04:05",
	})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

// Helper function to print or not
func (chroma *Chroma) println(a ...interface{}) {
	if !chroma.quiet {
		fmt.Fprintln(os.Stdout, a...)
	}
}

// Helper function to printf or not
func (chroma *Chroma) printf(format string, a ...interface{}) {
	if !chroma.quiet {
		fmt.Fprintf(os.Stdout, format, a...)
	}
}

// LogError is a helper function to trigger logging at the top level
func (chroma *Chroma) LogError(msg string, a ...interface{}) {
	logrus.Error(fmt.Sprintf(msg, a...))
}

// LogFatal exposes logging to caller for testing
func (chroma *Chroma) LogFatal(err error) {
	if chroma.debug {
		logrus.Fatalf("%+v", err)
	} else {
		logrus.Fatalf("%v", err)
	}
}
