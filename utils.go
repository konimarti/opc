package opc

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var logger *log.Logger

// Default is no logger
func init() {
	logger = newLogger(ioutil.Discard)
}

//Debug will set the logger to print to stderr
func Debug() {
	logger = newLogger(os.Stderr)
}

//SetLogWriter sets a user-defined writer for logger
func SetLogWriter(w io.Writer) {
	logger = newLogger(w)
}

//newLogger creats a log.Logger with standard settings
func newLogger(w io.Writer) *log.Logger {
	return log.New(w, "OPC ", log.LstdFlags)
}
