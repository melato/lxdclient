/*
Package lxdclient provides initialization code for connecting to an LXD server,
using the same configuration as used by the "lxc" command.

If using the unix socket, connect via the path:

$LXD_SOCKET, if defined

Otherwise, the first of these that exists:
  - $LXD_DIR/unix.socket
  - /var/lib/unix.socket
  - /var/snap/lxd/common/lxd/unix.socket

Otherwise, lxdclient looks for a config directory:

$LXD_CONF, if defined

Otherwise, the first of these that exists:
  - os.UserConfigDir()/lxc
  - $os.UserHomeDir()/snap/lxd/common/config

Where $var is the var environment variable.
*/
package lxdclient
