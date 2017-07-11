# Linen CNI plugin

A CNI plugin designed for overlay networks with [Open vSwitch](http://openvswitch.org).

## About Linen CNI plugin
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

## Example network configuration
Given the following network configurations for Node1(Master), Node2 and Node3:
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
	"ipMasq": true,
	"mtu": 1400,
	"hairpinMode": false,
	"vtepIPs": ["10.245.2.2", "10.245.2.3"],
	"ipam": {
		"type": "host-local",
		"subnet": "10.244.0.0/16",
		"rangeStart": "10.244.1.10",
			"rangeEnd": "10.244.1.150",
		"routes": [
			{ "dst": "0.0.0.0/0" }
		],
		"gateway": "10.244.1.1"
	}
}
EOF

$ tee /etc/cni/net.d/linen-cni.conf <<-'EOF'
{
	"name": "linen-demo-network",
	"type": "linen",
	"bridge": "kbr0",
	"ovsBridge": "br0",
	"isGateway": true,
	"isDefaultGateway": true,
	"forceAddress": false,
	"ipMasq": true,
	"mtu": 1400,
	"hairpinMode": false,
	"vtepIPs": ["10.245.2.2"],
	"ipam": {
		"type": "host-local",
		"subnet": "10.244.0.0/16",
		"rangeStart": "10.244.2.10",
			"rangeEnd": "10.244.2.150",
		"routes": [
			{ "dst": "0.0.0.0/0" }
		],
		"gateway": "10.244.2.1"
	}
}
EOF

$ tee /etc/cni/net.d/linen-cni.conf <<-'EOF'
{
	"name": "linen-demo-network",
	"type": "linen",
	"bridge": "kbr0",
	"ovsBridge": "br0",
	"isGateway": true,
	"isDefaultGateway": true,
	"forceAddress": false,
	"ipMasq": true,
	"mtu": 1400,
	"hairpinMode": false,
	"vtepIPs": ["10.245.2.2"],
	"ipam": {
		"type": "host-local",
		"subnet": "10.244.0.0/16",
		"rangeStart": "10.244.3.10",
			"rangeEnd": "10.244.3.150",
		"routes": [
			{ "dst": "0.0.0.0/0" }
		],
		"gateway": "10.244.3.1"
	}
}
EOF
```

### Network configuration reference

For **Linux Bridge plugin** options
- `name` (string, required): the name of the network.
- `type` (string, required): "bridge".
- `bridge` (string, optional): name of the bridge to use/create. Defaults to "cni0".
- `isGateway` (boolean, optional): assign an IP address to the bridge. Defaults to false.
- `isDefaultGateway` (boolean, optional): Sets isGateway to true and makes the assigned IP the default route. Defaults to false.
- `forceAddress` (boolean, optional): Indicates if a new IP address should be set if the previous value has been changed. Defaults to false.
- `ipMasq` (boolean, optional): set up IP Masquerade on the host for traffic originating from this network and destined outside of it. Defaults to false.
- `mtu` (integer, optional): explicitly set MTU to the specified value. Defaults to the value chosen by the kernel.
- `hairpinMode` (boolean, optional): set hairpin mode for interfaces on the bridge. Defaults to false.
- `ipam` (dictionary, required): IPAM configuration to be used for this network.

For **Open vSwitch Bridge plugin** options
- `ovsBridge`(string, required): name of the ovs bridge to use/create.
- `vtepIPs` (array, optional): array of the VXLAN tunnel end point IP addresses

## Usage in Kubernetes
1. Create Linen CNI configuration file in the `/etc/cni/net.d/linen-cni.conf` directories.
2. Make sure that the linen binary are in the `/opt/cni/bin` directories directories.
3. Test to create a POD/Deployment