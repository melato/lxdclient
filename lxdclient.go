package lxdclient

import (
	"fmt"

	lxd "github.com/canonical/lxd/client"
	"github.com/canonical/lxd/lxc/config"
	"github.com/canonical/lxd/shared"
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

	conf       *config.Config
	rootServer lxd.InstanceServer
}

func (t *LxdClient) Init() error {
	return nil
}

func connectUnix() (lxd.InstanceServer, error) {
	unixSocket, err := UnixSocket()
	if err != nil {
		return nil, err
	}
	if unixSocket == "" || !shared.PathExists(unixSocket) {
		return nil, fmt.Errorf("no such unix socket: %s", unixSocket)
	}
	if Trace {
		fmt.Printf("using unix socket: %s\n", unixSocket)
	}
	server, err := lxd.ConnectLXDUnix(unixSocket, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", unixSocket, err)
	}
	return server, nil
}

func (t *LxdClient) Configured() error {
	if t.ForceLocal {
		server, err := connectUnix()
		if err != nil {
			return err
		}
		t.rootServer = server
		t.conf = config.NewConfig("", true)
	} else {
		conf, err := LoadConfig()
		if err != nil {
			return err
		}
		t.conf = conf
	}
	return nil
}

// RootServer - return the unqualified (no project) LXD instance server
func (t *LxdClient) RootServer() (lxd.InstanceServer, error) {
	if t.rootServer == nil {
		remote, ok := t.conf.Remotes[t.conf.DefaultRemote]
		if ok && remote.Addr == "unix://" {
			server, err := connectUnix()
			if err != nil {
				return nil, err
			}
			t.rootServer = server
		} else {
			d, err := t.conf.GetInstanceServer(t.conf.DefaultRemote)
			if err != nil {
				return nil, err
			}
			t.rootServer = d
		}
	}
	return t.rootServer, nil
}

// RootServer - return the LXD instance server for the specified project
// If project is empty, use the default project
func (t *LxdClient) ProjectServer(project string) (lxd.InstanceServer, error) {
	if project == "" {
		project = t.CurrentProject()
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

func (t *LxdClient) CurrentProject() string {
	if t.Project != "" {
		return t.Project
	}
	if t.conf != nil {
		remote, exists := t.conf.Remotes[t.conf.DefaultRemote]
		if exists {
			return remote.Project
		}
	}
	return ""
}

func (t *LxdClient) Config() *config.Config {
	return t.conf
}
