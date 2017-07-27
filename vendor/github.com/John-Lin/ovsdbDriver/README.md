ovsdbDriver
====
A libovsdb wrapper for operating [Open vSwitch](http://openvswitch.org/) via Go.

This library is a fork of the ovsdbDriver functionality in [contiv/ofnet](https://github.com/contiv/ofnet).

## Install 

```
$ go get -u github.com/John-Lin/ovsdbDriver
```

## Usage with ovsdb-server

You might need to create OVS bridge on given IP and TCP port, please make sure that you have set OVSDB listener.

```
$ ovs-vsctl set-manager ptcp:6640
```

Create bridge should assign IP and TCP port.

```go
ovsDriver = ovsdbDriver.NewOvsDriver("ovsbr", "127.0.0.1", 6640)
```

Otherwise, `ovsdb-server` connects to the Unix domain server socket and the default path is `unix:/var/run/openvswitch/db.sock` .

```go
ovsDriver = ovsDriver = NewOvsDriverWithUnix("br0")
```

## Example
```go
package main

import "github.com/John-Lin/ovsdbDriver"

var ovsDriver *ovsdbDriver.OvsDriver

func main() {
    // Create an OVS bridge to the Unix domain server socket.
    ovsDriver = NewOvsDriverWithUnix("br0")
    
    // Create an OVS bridge to the given IP and TCP port.
    // ovsDriver = ovsdbDriver.NewOvsDriver("ovsbr", "127.0.0.1", 6640)
    
    // Add ovsbr as a internal port without vlan tag (0)
    ovsDriver.CreatePort("ovsbr", "internal", 0)
}
```

Use `ovs-vsctl show` to check bridge information.

```
root@dev:~# ovs-vsctl show
e650c132-8c99-44b4-aa50-c640645f4f18
    Manager "ptcp:6640"
    Bridge "ovsbr"
        Port "ovsbr"
            Interface "ovsbr"
                type: internal
    ovs_version: "2.5.2"
```

## Related
- [libovsdb](https://github.com/socketplane/libovsdb) is an OVSDB library which is originally developed by SocketPlane.


