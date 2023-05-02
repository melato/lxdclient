package lxdclient

import (
	"os"
	"path/filepath"

	"github.com/lxc/lxd/shared"
)

/*
UnixSocket returns the default LXD Unix Socket
See https://github.com/lxc/lxd/blob/master/client/connection.go
We use a similar strategy as connection.go ConnectLXDUnixWithContext(),
except that we add /var/snap/lxd/common/lxd/ to the LXD directories searched.
If no LXD directory is found, return "", nil.
*/
func UnixSocket() (string, error) {
	path := os.Getenv("LXD_SOCKET")
	if path != "" {
		return path, nil
	}
	lxdDir := os.Getenv("LXD_DIR")
	if lxdDir == "" {
		lxdDir = "/var/lib/lxd"
	}
	if !shared.PathExists(lxdDir) {
		lxdDir = "/var/snap/lxd/common/lxd"
	}
	if !shared.PathExists(lxdDir) {
		return "", nil
	}

	path = filepath.Join(lxdDir, "unix.socket")
	return path, nil
}

/*
ConfigDir returns the default LXD client configuration directory.
See https://github.com/lxc/lxd/blob/master/client/main.go
We use a similar strategy as main.go func (c *cmdGlobal) PreRun,
except that we add $HOME/snap/lxd/common/config/ to the LXD config directories searched.
If no LXD CONF directory is found, return "", nil.
*/
func ConfigDir() (string, error) {
	configDir := os.Getenv("LXD_CONF")
	if configDir != "" {
		return configDir, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	userConfigDir, err := os.UserConfigDir()
	if err == nil && userConfigDir != "" {
		configDir = filepath.Join(userConfigDir, "lxc")
		if shared.PathExists(configDir) {
			return configDir, nil
		}
	}
	configDir = filepath.Join(home, "snap", "lxd", "common", "config")
	if shared.PathExists(configDir) {
		return configDir, nil
	}
	return "", nil
}