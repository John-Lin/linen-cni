ovsdb
====
A libovsdb wrapper for operating [Open vSwitch](http://openvswitch.org/) via Go.

This library is a fork of the ovsdbDriver functionality in [contiv/ofnet](https://github.com/contiv/ofnet) and makes modifications for supporting Unix domain socket to connect OVSDB.

## Install 

```
$ go get -u github.com/John-Lin/ovsdb
```

## Usage with ovsdb-server

You might need to create OVS bridge on given IP and TCP port, please make sure that you have set OVSDB listener.

```
$ ovs-vsctl set-manager ptcp:6640
```

Create bridge should assign IP and TCP port.

```go
ovsDriver = ovsdb.NewOvsDriver("ovsbr", "127.0.0.1", 6640)
```

Otherwise, `ovsdb-server` connects to the Unix domain server socket and the default path is `unix:/var/run/openvswitch/db.sock` .

```go
ovsDriver = ovsdb.NewOvsDriverWithUnix("br0")
```

## Example
```go
package main

import "github.com/John-Lin/ovsdb"

var ovsDriver *ovsdb.OvsDriver

func main() {
    // Create an OVS bridge to the Unix domain server socket.
    ovsDriver = ovsdb.NewOvsDriverWithUnix("br0")

    // Create an OVS bridge to the given IP and TCP port.
    // ovsDriver = ovsdb.NewOvsDriver("br0", "127.0.0.1", 6640)

    // Add br0 as a internal port without vlan tag (0)
    ovsDriver.CreatePort("br0", "internal", 0)
}
```

Use `ovs-vsctl show` to check bridge information.

```
root@dev:~# ovs-vsctl show
82040598-7050-4320-b946-1d4380fabc73
    Bridge "br0"
        Port "br0"
            Interface "br0"
                type: internal
    ovs_version: "2.5.2"
```

## Related
- [libovsdb](https://github.com/socketplane/libovsdb) is an OVSDB library which is originally developed by SocketPlane.
- [contiv/ofnet](https://github.com/contiv/ofnet) is openflow networking library.


