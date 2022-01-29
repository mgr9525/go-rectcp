package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	println("hello rectcp:", Version)
	regs()
	kingpin.Version(Version)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
