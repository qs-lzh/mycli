package main

import (
	"log"

	"github.com/qs-lzh/mycli/commands"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	commands.Execute()
}
