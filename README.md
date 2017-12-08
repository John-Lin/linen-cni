# Linen CNI plugin

A CNI plugin designed for overlay networks with [Open vSwitch](http://openvswitch.org).

## About Linen CNI plugin
Linen provides a convenient way to easily setup networking between pods across nodes. To support multi-host overlay networking and large scale isolatio, VxLAN tunnel end point (VTEP) is used instead of GRE. Linen creates an OVS bridge and added as a port to the linux bridge.

This CNI plugin implementation was inspired by the document from [Kubernetes OVS networking](https://kubernetes.io/docs/admin/ovs-networking/) and designed to meet the requirements of SDN environment.

Please read [CNI](https://github.com/containernetworking/cni/blob/master/SPEC.md) for more detail on container networking.

## Prerequisite
```
$ sudo apt-get install openvswitch-switch
```

# Should I use this or ovn-kubernetes?
ovn-kubernetes provides more advanced features and use vRouter (Layer 3 approach) to achieve multi-host networking. 
If you're going to create vRouters and vSwitches to build any network topologies you desire, ovn-kubernetes is a complete solution. 

This CNI plugin creates only vSwitches in each node and uses VxLAN for achieving network overlay.
For the PODs in cluster are managed by linux bridges and the IP allocation is configured through `IPAM` plugin.

# Kubernetes
**Kubernetes 1.7+ and CNI 0.6.0 are required**.

Linen CNI is not only a plugin which support for network namespace (e.g., docker, ip-netns), but also a option for Kubernetes cluster networking.

## Usage
1. Create a Linen CNI configuration list file in the `/etc/cni/net.d/linen.conflist` directories.
2. Make sure that the `linen`, `bridge` and `host-local` binaries are in the `/opt/cni/bin` directories directories.
3. (Optional) Create a daemon set to manager ovsdb `kubectl create -f flaxd.yaml`.
3. Test to create a POD/Deployment.

## Architecture

### Management Workflow

- `flax daemon`: Runs on each host in order to monitor new node join and add it to current overlay network.
- `linen-cni`: Executed by the container runtime and set up the network stack for containers.

<p align="center">
    <img src="/images/mgmt-workflow.png" width="541" />
</p>

### Packet Processing

To provide overlay network, Linen utilize Open vSwitch to create VxLAN tunneling in the backend.

<p align="center">
    <img src="/images/ovs-networking.png" width="586" />
</p>

## Example network configuration
Please check example network configuration in the `examples` folder


### Network configuration reference

Linen is a chained plugin. It always comes after `bridge` plugin, so configure Linux Bridge is needed.

For the **Linux Bridge plugin** options
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
- `promiscMode` (boolean, optional): set promiscuous mode on the bridge. Defaults to false.

For the **Open vSwitch Bridge plugin** options
- `isMaster`(boolean, optional): sets isMaster to true if the host is the Kubernetes master node in cluster. Defaults to false.
- `ovsBridge`(string, optional): name of the ovs bridge to use/create.
- `vtepIPs` (list, optional): list of the VxLAN tunnel end point IP addresses.
- `controller` (string, optional): sets SDN controller, assigns an IP address, port number like `192.168.100.20:6653`. Controller is not not essential for overlay network. 

## Build
You may need to build the binary from source. The "build-essential" package is required.

```
$ sudo apt-get install build-essential
```

Execute `build.sh`

```
$ ./build.sh
```

When build succeed, binary will be in the `bin` folder.

# Example
Linen-CNI also provides a vagrantfile to help you setup a demo environment to try Linen-CNI.

## Environment
You should install vagrant in your system and make sure everything goes well.

## Setup Linen-CNI
- Change directory to `Linen-CNI` and type `vagrant up` to init a virtual machine.
- Use ssh to connect vagrant VM via `vagrant ssh`.
- Type following commang to build the `linen-cni` binary and move it to CNI directory.
```
cd linen-cni
sh build.sh
cp bin/linen ../cni/ 
```
- We need to provide a CNI config for `Linen-CNI`, and you can use build-in config from example directory. Use following command to copy the config to `/root` directory.

```
sudo cp examples/master.linen.conflist  /root/linen.conflist
```

## Create NS
In this vagrant environment, we don't install docker related services but you can use `namespace(ns)` to test `Linen-CNI`.
Type following command to create a namespace named ns1

```
sudo ip netns add ns1
```

## Start CNI
We have setup Linen-CNI environement and testing namespace(ns1), we can use the following commands to inform CNI to add a network for the namespace.

```
cd ~/cni
sudo CNI_PATH=`pwd` NETCONFPATH=/root ./cnitool \ add linen-network /var/run/netns/ns1
```
and the result looks like below
```
{
    "cniVersion": "0.3.1",
    "interfaces": [
        {
            "name": "veth7df4d2c0",
            "mac": "56:b1:e8:32:e4:b7"
        },
        {
            "name": "eth0",
            "mac": "0a:58:0a:f4:01:0a",
            "sandbox": "/var/run/netns/ns1"
        }
    ],
    "ips": [
        {
            "version": "4",
            "interface": 2,
            "address": "10.244.1.10/16",
            "gateway": "10.244.1.1"
        }
    ],
    "routes": [
        {
            "dst": "0.0.0.0/0"
        },
        {
            "dst": "0.0.0.0/0",
            "gw": "10.244.1.1"
        }
    ],
    "dns": {}
}
```

Now, we can use some tools to help us check the current network setting, for example.  
You can use `ovs-vsctl show` to show current OVS setting and it looks like: 

```
e6289dc2-a181-4316-b902-a50fc6d854b6
    Bridge "br0"
        Controller "tcp:192.168.2.100:6653"
        fail_mode: standalone
        Port "vxif10_245_2_2"
            Interface "vxif10_245_2_2"
                type: vxlan
                options: {key=flow, remote_ip="10.245.2.2"}
        Port "br0"
            Interface "br0"
                type: internal
        Port "vxif10_245_2_3"
            Interface "vxif10_245_2_3"
                type: vxlan
                options: {key=flow, remote_ip="10.245.2.3"}
    ovs_version: "2.5.2"
```
In this setting, the OVS will try to connect to Openflow controller (it not exist, change to L2 bridge mode) and it also contains three ports, including two vxlan ports.  

Besides, you can use `brctl show` to see that the OVS bridge (br0) has been attached to Linux bridge(kbr).

```
bridge name     bridge id               STP enabled     interfaces
kbr0            8000.0a580af40101       no              br0
                                                        veth7df4d2c0
```

If you want to check the namepsace's networking setting, you can use `sudo ip netns exec ns1 ifconfig` to see it's IP config.
```
ubuntu@dev:~$ sudo ip netns exec ns1 ifconfig
eth0      Link encap:Ethernet  HWaddr 0a:58:0a:f4:01:0a
          inet addr:10.244.1.10  Bcast:0.0.0.0  Mask:255.255.0.0
          inet6 addr: fe80::bc15:faff:fe6b:b414/64 Scope:Link
          UP BROADCAST RUNNING MULTICAST  MTU:1400  Metric:1
          RX packets:18 errors:0 dropped:0 overruns:0 frame:0
          TX packets:10 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:0
          RX bytes:1476 (1.4 KB)  TX bytes:828 (828.0 B)
```
