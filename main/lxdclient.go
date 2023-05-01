package main

import (
	_ "embed"
	"fmt"

	"melato.org/command"
	"melato.org/command/usage"
	"melato.org/lxdclient"
)

//go:embed usage.yaml
var usageData []byte

type App struct {
	lxdclient.LxdClient
}

func (t *App) ListContainers() error {
	server, err := t.CurrentServer()
	if err != nil {
		return err
	}
	names, err := server.GetContainerNames()
	if err != nil {
		return err
	}
	for _, name := range names {
		fmt.Printf("%s\n", name)
	}
	return nil
}

func main() {
	cmd := &command.SimpleCommand{}
	var app App
	cmd.Flags(&app).Command("list").RunFunc(app.ListContainers)

	usage.Apply(cmd, usageData)
	command.Main(cmd)
}
