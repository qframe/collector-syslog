package main

import (
	"log"
	"github.com/zpatrick/go-config"
	"github.com/qnib/qframe-types"
	"os"
	"github.com/qframe/collector-syslog"
)

func main() {
	qChan := qtypes.NewQChan()
	qChan.Broadcast()
	if len(os.Args) != 2 {
		log.Fatal("usage: ./syslog <port>")

	}
	port := os.Args[1]
	cfgMap := map[string]string{
		"collector.syslog.port": port,
	}
	cfg := config.NewConfig([]config.Provider{config.NewStatic(cfgMap)})

	p, err := qcollector_syslog.New(qChan, cfg, "syslog")
	if err != nil {
		log.Fatalf("[EE] Failed to create collector: %v", err)
	}
	go p.Run()
	bg := p.QChan.Data.Join()
	for {
		val := <- bg.Read
		log.Printf("%v", val)
	}
}
