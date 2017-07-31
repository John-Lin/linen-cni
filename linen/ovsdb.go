package main

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/John-Lin/ovsdbDriver"
	"github.com/vishvananda/netlink"
)

var ovsDriver *ovsdbDriver.OvsDriver

func ensureOVSBridge(OVSBrName string) (*netlink.Bridge, error) {

	// create a ovs bridge
	ovsDriver = ovsdbDriver.NewOvsDriverWithUnix(OVSBrName)

	// Create an internal port in OVS
	ovsDriver.CreatePort(OVSBrName, "internal", 0)
	// err := ovsDriver.CreatePort(OVSBrName, "internal", 0)
	// if err != nil {
	// 	return fmt.Errorf("Error creating the port. Err: %v", err)
	// }

	time.Sleep(300 * time.Millisecond)

	// finds a link by name and returns a pointer to the object.
	// ovsbr, _ := netlink.LinkByName(OVSBrName)
	ovsbrLink, err := netlink.LinkByName(OVSBrName)
	if err != nil {
		return nil, fmt.Errorf("could not lookup link on ensureOVSBridge %q: %v", OVSBrName, err)
	}

	// enables the link device
	if err := netlink.LinkSetUp(ovsbrLink); err != nil {
		return nil, err
	}

	ovsbr, _ := bridgeByName(OVSBrName)

	return ovsbr, nil
}

func setupCtrlerToOVS(config *LinenConf) error {
	// setup SDN controller for ovs bridge
	host, port, err := net.SplitHostPort(config.RuntimeConfig.OVS.Controller)
	if err != nil {
		return fmt.Errorf("Invalid controller IP and port. Err: %v", err)
	}

	uPort, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		return fmt.Errorf("Invalid controller port number. Err: %v", err)
	}

	err = ovsDriver.AddController(host, uint16(uPort))
	if err != nil {
		return fmt.Errorf("Error adding controller to OVS. Err: %v", err)
	}
	return nil
}

func setupOVSBridge(config *LinenConf) (*netlink.Bridge, error) {
	// create ovs bridge
	ovsbr, err := ensureOVSBridge(config.RuntimeConfig.OVS.OVSBrName)
	if err != nil {
		return nil, err
	}

	return ovsbr, nil
}

func setupVTEPs(config *LinenConf) error {
	for i := 0; i < len(config.RuntimeConfig.OVS.VtepIPs); i++ {

		// Create interface name for VTEP
		intfName := vxlanIfName(config.RuntimeConfig.OVS.VtepIPs[i])

		// Check if it already exists
		isPresent, vsifName := ovsDriver.IsVtepPresent(config.RuntimeConfig.OVS.VtepIPs[i])
		if !isPresent || (vsifName != intfName) {
			// create VTEP
			err := ovsDriver.CreateVtep(intfName, config.RuntimeConfig.OVS.VtepIPs[i])
			if err != nil {
				return fmt.Errorf("Error creating VTEP port %s. Err: %v", intfName, err)
			}
		}

	}

	return nil
}

func addOVSBridgeToBridge(config *LinenConf) error {
	ovsbrLink, err := netlink.LinkByName(config.RuntimeConfig.OVS.OVSBrName)
	if err != nil {
		return fmt.Errorf("could not lookup link on addOVSBridgeToBridge %q: %v", config.RuntimeConfig.OVS.OVSBrName, err)
	}
	// The first element for Interfaces is brInterface
	br, err := bridgeByName(config.PrevResult.Interfaces[0].Name)
	if err != nil {
		return err
	}

	// Adding the interface into the bridge is done by setting its master to bridge_name()
	// netlink.LinkSetMaster(ovsbrLink, br)
	if err := netlink.LinkSetMaster(ovsbrLink, br); err != nil {
		return fmt.Errorf("failed to LinkSetMaster %v", err)
	}
	return nil
}

func bridgeByName(name string) (*netlink.Bridge, error) {
	l, err := netlink.LinkByName(name)
	if err != nil {
		return nil, fmt.Errorf("could not lookup %q: %v", name, err)
	}
	br, ok := l.(*netlink.Bridge)
	if !ok {
		return nil, fmt.Errorf("%q already exists but is not a bridge", name)
	}
	return br, nil
}
