package qcollector_syslog

import (
	"github.com/zpatrick/go-config"
	"gopkg.in/mcuadros/go-syslog.v2"
	"fmt"
	"github.com/qnib/qframe-types"
	"github.com/qframe/types/messages"
	"time"
)

const (
	version = "0.1.0"
	pluginTyp = "collector"
	pluginPkg = "syslog"
)

type Plugin struct {
	qtypes.Plugin
	slogChannel syslog.LogPartsChannel
	sPort string
}

func New(qChan qtypes.QChan, cfg *config.Config, name string) (Plugin, error) {
	return Plugin{
		Plugin: qtypes.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg, name, version),
		slogChannel: make(syslog.LogPartsChannel),
	}, nil
}

func (p *Plugin) Run() {
	p.Log("notice", fmt.Sprintf("Start collector v%s", version))
	var err error
	p.sPort, err = p.CfgString("port")
	if err != nil {
		p.Log("error", "No configuration for port found")
		return
	}
	addr := fmt.Sprintf("0.1.0.0:%s", p.sPort)
	handler := syslog.NewChannelHandler(p.slogChannel)
	server := syslog.NewServer()
	server.SetFormat(syslog.RFC5424)
	server.SetHandler(handler)
	server.ListenUDP(addr)
	server.ListenTCP(addr)
	server.Boot()
	go p.PushMessages()
	server.Wait()
}

func (p *Plugin) PushMessages() {
	layout := "2006-01-02 15:04:05 -0700 MST"
	for logParts := range p.slogChannel {
		t, err := time.Parse(layout, fmt.Sprintf("%s", logParts["timestamp"]))
		if err != nil {
			p.Log("error", fmt.Sprintf("Could not parse time from: %s", logParts["timestamp"]))
			t = time.Now()
		}
		b := qtypes_messages.NewTimedBase(p.Name, t)
		qm := qtypes_messages.NewMessage(b, fmt.Sprintf("%s", logParts["message"]))
		qm.Tags["msg_id"] = fmt.Sprintf("%s", logParts["msg_id"])
		qm.Tags["priority"] = fmt.Sprintf("%d", logParts["priority"])
		qm.Tags["severity"] = fmt.Sprintf("%d", logParts["severity"])
		qm.Tags["app_name"] = fmt.Sprintf("%s", logParts["app_name"])
		qm.Tags["facility"] = fmt.Sprintf("%d", logParts["facility"])
		qm.Tags["hostname"] = fmt.Sprintf("%s", logParts["hostname"])
		p.QChan.Data.Send(qm)
	}
}

