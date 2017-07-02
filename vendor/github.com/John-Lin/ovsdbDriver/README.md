ovsdbDriver
====
A libovsdb wrapper for operating [Open vSwitch](http://openvswitch.org/) via Go.

This library is a fork of the ovsdbDriver functionality in [contiv/ofnet](https://github.com/contiv/ofnet).

## Install 

```
$ go get github.com/John-Lin/ovsdbDriver
```

## Example
```go
package main
import "github.com/John-Lin/ovsdbDriver"

var ovsDriver *ovsdbDriver.OvsDriver

func main() {
    // Create an OVS bridge
    ovsDriver = NewOvsDriver("ovsbr")
    // Add ovsbr as a internal port without vlan tag (0)
    ovsDriver.CreatePort("ovsbr", "internal", 0)
}
```

Use `ovs-vsctl show` to check bridge information.

```
root@dev:~# ovs-vsctl show
e650c132-8c99-44b4-aa50-c640645f4f18
    Bridge "ovsbr"
        Port "ovsbr"
            Interface "ovsbr"
                type: internal
    ovs_version: "2.5.2"
```

## Related
- [libovsdb](https://github.com/socketplane/libovsdb) is an OVSDB library which is originally developed by SocketPlane.


