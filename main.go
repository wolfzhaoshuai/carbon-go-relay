package main

import (
	"flag"
	"fmt"
	"os"

	"carbon-go-relay/apm"
	"carbon-go-relay/global"
	"carbon-go-relay/http"
	"carbon-go-relay/receiver"
	"carbon-go-relay/sender"
	"carbon-go-relay/utils"
)

func main() {
	utils.Zlog.Info("carbon-c-relay carbon-go-relay start")
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println(global.VERSION)
		os.Exit(0)
	}
	// global config
	global.ParseConfig(*cfg)
	global.InitalizeGlobalConstants()

	http.Start()
	receiver.Start()
	sender.Start()
	apm.Start()

	select {}
}
