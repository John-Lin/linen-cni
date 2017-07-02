# Linen CNI plugin

A CNI plugin designed for overlay networks with Open vSwitch.

# Linen Configuration file

``
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
