package chroma

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// Setup log formatting
func (chroma *CHROMA) setupLogging() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "02-01-2006 15:04:05",
	})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

// Helper function to print or not
func (chroma *CHROMA) println(a ...interface{}) {
	if !chroma.quiet {
		fmt.Fprintln(os.Stdout, a...)
	}
}

// Helper function to printf or not
func (chroma *CHROMA) printf(format string, a ...interface{}) {
	if !chroma.quiet {
		fmt.Fprintf(os.Stdout, format, a...)
	}
}

// Helper function to choose between detailed or not
func (chroma *CHROMA) logFatal(err error) {
	if chroma.debug {
		logrus.Fatalf("%+v", err)
	} else {
		logrus.Fatalf("%v", err)
	}
}
