//+build profile

package main

import (
	// web profile
	_ "net/http/pprof"
)

func init() {
}
