package opc

import (
	"io/ioutil"
	"log"
	"os"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func Debug() {
	log.SetOutput(os.Stderr)
}
