package main

import (
	_ "embed"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
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

func (t *App) ListProfiles() error {
	server, err := t.CurrentServer()
	if err != nil {
		return err
	}
	profiles, err := server.GetProfiles()
	if err != nil {
		return err
	}
	for _, p := range profiles {
		fmt.Printf("%s\n", p.Name)
	}
	return nil
}

func (t *App) PrintConfig() error {
	conf := t.Config()
	data, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	os.Stdout.Write(data)
	fmt.Println()
	return nil
}

func main() {
	cmd := &command.SimpleCommand{}
	var app App
	cmd.Flags(&app).Command("conf").RunFunc(app.PrintConfig)
	cmd.Flags(&app).Command("containers").RunFunc(app.ListContainers)
	cmd.Flags(&app).Command("profiles").RunFunc(app.ListProfiles)

	usage.Apply(cmd, usageData)
	command.Main(cmd)
}
