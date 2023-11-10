package core

import (
	"fmt"
	"github.com/projectdiscovery/goflags"
	"syscall"
)

type Options struct {
	Webhook   string
	WildCards goflags.StringSlice
	Delay     int
	Vdp       bool
	Log       bool
}

// NewParser function creates and returns a new instance of the Options struct.
func NewParser() *Options {
	return &Options{}
}

// Parse function parses command-line arguments, sets options, and displays a banner.
func (o *Options) Parse() {
	flagSet := goflags.NewFlagSet()
	flagSet.StringSliceVarP(&o.WildCards, "domains", "d", nil, "domain of targets", goflags.CommaSeparatedStringSliceOptions)
	flagSet.SetDescription("ScopeDetective Program To Get Latest Scope In HackerOne")
	flagSet.StringVar(&o.Webhook, "webhook", "", "discord webhook url")
	flagSet.IntVar(&o.Delay, "delay", 10, "delay (min, default 10)")
	flagSet.BoolVar(&o.Vdp, "vdp", false, "get vdp program")
	flagSet.BoolVar(&o.Log, "log", false, "send log")
	_ = flagSet.Parse()

	showBanner()

	if o.Webhook == "" {
		fmt.Println("\033[31m[!] Usage: ScopeDetective -webhook <webhook> -delay <delay> \033[0m")
		syscall.Exit(0)
	}
}
