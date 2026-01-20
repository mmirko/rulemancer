package main

/*
#cgo CFLAGS: -I core
#cgo LDFLAGS: -L core -lclips -lm
#include <stdlib.h>
*/
import "C"

import (
	"github.com/mmirko/rulemancer/cmd"
)

func main() {
	cmd.Execute()
}
