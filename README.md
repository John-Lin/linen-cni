# Linen CNI plugin

A CNI plugin designed for overlay networks with [Open vSwitch](http://openvswitch.org).

# About Linen CNI plugin
Linen provides a convenient way to easily setup networking between pods across nodes. To support multi-host overlay networking and large scale isolatio, VXLAN tunnel end point (VTEP) is used instead of GRE. Linen creates an OVS bridge and added as a port to the linux bridge.

This CNI plugin implementation was inspired by the document from [Kubernetes OVS networking](https://kubernetes.io/docs/admin/ovs-networking/) and designed to meet the requirements of SDN environment.

Please read [CNI](https://github.com/containernetworking/cni/blob/master/SPEC.md) for more detail on container networking.

## Architecture

![OVS Networking](/images/ovs-networking.png)

## Build

```
$ ./build.sh
```

when build succeed binary will be in the `bin` folder.

## Linen Network Configuration
Given the following network configuration:
```
$ tee /etc/cni/net.d/linen-cni.conf <<-'EOF'
{
	"name": "linen-demo-network",
	"type": "linen",
	"bridge": "kbr0",
	"ovsBridge": "br0",
        "isGateway": true,
	"isDefaultGateway": true,
	"forceAddress": false,
	"ipMasq": false,
        "mtu": 1400,
	"hairpinMode": false,
	"vtepIP": ["192.168.120.10", "192.168.60.5", "192.168.30.1"],
	"ipam": {
		"type": "host-local",
		"subnet": "10.244.0.0/16",
                "routes": [
                        { "dst": "0.0.0.0/0" }
                ],
		"gateway": "10.244.1.1"
	}
}
EOF
```

## Usage in Kubernetes
1. Create Linen CNI configuration file in the `/etc/cni/net.d/linen-cni.conf` directories.
2. Make sure that the linen binary are in the `/opt/cni/bin` directories directories.
3. Test to create a POD/Deployment