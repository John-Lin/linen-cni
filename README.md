# Linen CNI plugin

A CNI plugin designed for overlay networks with Open vSwitch.

## Network Architecture
Coming soon

## Build

```
$ ./build.sh
```

when build succeed binary will be in the `bin` folder.

## Linen Configuration file
Here is an example for create an overlay network using OVS
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
