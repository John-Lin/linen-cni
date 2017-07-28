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
- `promiscMode` (boolean, optional): set promiscuous mode on the bridge. Defaults to false.

For **Open vSwitch Bridge plugin** options
- `isMaster`(boolean, optional): sets isMaster to true if the host is the Kubernetes master node in cluster. Defaults to false.
- `bridge` (string, optional): name of the bridge to connect to ovs bridge. Defaults to "cni0".
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
