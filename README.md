This is a small Go library that connects to an LXD server,
with the goal of using the same configuration used by the "lxc" command.

It Uses code from https://github.com/lxc/lxd

The lxc code use these default paths:
- /var/lib/lxd/unix.socket (client/connection.go)
- $HOME/.config/lxc (lxc/main.go)

But when LXD is installed as a snap, "lxc" seems to use the following:

- /var/snap/lxd/common/lxd/unix.socket 
- $HOME/snap/lxd/common/config/


lxdclient uses all of the above, whichever it finds first.

See go doc for the full list.