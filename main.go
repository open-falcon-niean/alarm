package main

import (
	"flag"
	"fmt"
	"github.com/open-falcon/alarm/cron"
	"github.com/open-falcon/alarm/exco"
	"github.com/open-falcon/alarm/g"
	"github.com/open-falcon/alarm/http"
	"os"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	g.ParseConfig(*cfg)
	g.InitRedisConnPool()

	// exco
	exco.Start()

	go http.Start()

	go cron.ReadHighEvent()
	go cron.ReadLowEvent()
	go cron.CombineSms()
	go cron.CombineMail()

	select {}
}
