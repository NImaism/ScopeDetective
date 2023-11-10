package main

import (
	"fmt"
	"github.com/NImaism/ScopeDetective/core"
	"os"
	"os/signal"
)

// main function sets up the program, handles interruptions, and runs the system.
func main() {
	options := core.NewParser()
	options.Parse()

	system := core.New(core.NewMessager(options), *options)
	fresh := core.NewFresh(core.NewMessager(options), *options)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println("\n\033[35m[+] Bye See You Later\033[0m")
			os.Exit(1)
		}
	}()

	go fresh.Run()
	system.Run()
}
