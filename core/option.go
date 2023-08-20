package core

import (
	"fmt"
	"github.com/projectdiscovery/goflags"
	"syscall"
)

type Options struct {
	Webhook string
	Delay   int
}

// NewParser function creates and returns a new instance of the Options struct.
func NewParser() *Options {
	return &Options{}
}

// Parse function parses command-line arguments, sets options, and displays a banner.
func (o *Options) Parse() {
	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription("ScopeDetective Program To Get Latest Scope In HackerOne")
	flagSet.StringVar(&o.Webhook, "webhook", "", "discord webhook url")
	flagSet.IntVar(&o.Delay, "delay", 10, "delay (min, default 10)")

	_ = flagSet.Parse()

	showBanner()

	if o.Webhook == "" {
		fmt.Println("\033[31m[!] Usage: ScopeDetective -webhook <webhook> -delay <delay> \033[0m")
		syscall.Exit(0)
	}
}
