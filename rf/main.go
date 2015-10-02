package main

import (
	"flag"

	"github.com/nextrevision/runfile"
)

func main() {
	flag.Parse()
	rf.Run(flag.Args())
}
