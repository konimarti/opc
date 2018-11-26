package opc

import (
	"io/ioutil"
	"log"
	"os"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

//Debug will print out more information about the package.
func Debug() {
	log.SetOutput(os.Stderr)
}
