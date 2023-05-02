package lxdclient

import (
	"fmt"
	"os"
	"path/filepath"

	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/lxc/config"
	"github.com/lxc/lxd/shared"
)

var Trace bool

/*
LxdClient provides an lxd.InstanceServer
Use it by calling these methods, in order:
Configured(), CurrentServer().
*/
type LxdClient struct {
	ForceLocal bool   `name:"force-local" usage:"Force using the local unix socket"`
	Project    string `name:"project" usage:"Override the default project"`

	confPath   string
	conf       *config.Config
	rootServer lxd.InstanceServer
}

func (c *LxdClient) Configured() error {
	if c.ForceLocal {
		unixSocket, err := UnixSocket()
		if err != nil {
			return err
		}
		if unixSocket == "" || !shared.PathExists(unixSocket) {
			return fmt.Errorf("no such unix socket: %s", unixSocket)
		}
		if Trace {
			fmt.Printf("using unix socket: %s\n", unixSocket)
		}
		server, err := lxd.ConnectLXDUnix(unixSocket, nil)
		if err != nil {
			return fmt.Errorf("%s: %w", unixSocket, err)
		}
		c.rootServer = server
		c.conf = config.NewConfig("", true)
	} else {
		configDir, err := ConfigDir()
		if err != nil {
			return err
		}
		if configDir == "" || !shared.PathExists(configDir) {
			return fmt.Errorf("config dir not found: %s", configDir)
		}
		c.confPath = os.ExpandEnv(filepath.Join(configDir, "config.yml"))
		c.conf, err = config.LoadConfig(c.confPath)
		if err != nil {
			return err
		}
	}
	return nil
}

// RootServer - return the unqualified (no project) LXD instance server
func (t *LxdClient) RootServer() (lxd.InstanceServer, error) {
	if t.rootServer == nil {
		d, err := t.conf.GetInstanceServer(t.conf.DefaultRemote)
		if err != nil {
			return nil, err
		}
		t.rootServer = d
	}
	return t.rootServer, nil
}

// RootServer - return the LXD instance server for the specified project
// If project is empty, use the default project
func (t *LxdClient) ProjectServer(project string) (lxd.InstanceServer, error) {
	var err error
	if project == "" {
		project = t.Project
	}
	if project == "" {
		remote, ok := t.conf.Remotes[t.conf.DefaultRemote]
		if ok {
			project = remote.Project
		}
	}
	if project == "" {
		project = "default"
	}
	server, err := t.RootServer()
	if err != nil {
		return nil, err
	}
	return server.UseProject(project), nil
}

// RootServer - return the LXD instance server for the current project
func (t *LxdClient) CurrentServer() (lxd.InstanceServer, error) {
	return t.ProjectServer("")
}

// Config - return the LXD *Config
func (c *LxdClient) Config() *config.Config {
	return c.conf
}
