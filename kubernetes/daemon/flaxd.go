package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/John-Lin/ovsdbDriver"
	log "github.com/Sirupsen/logrus"
	"github.com/containernetworking/cni/pkg/types"
)

var ovsDriver *ovsdbDriver.OvsDriver

type OVS struct {
	IsMaster  bool     `json:"isMaster"`
	OVSBrName string   `json:"ovsBridge"`
	VtepIPs   []string `json:"vtepIPs"`
}

type NetConf struct {
	types.NetConf
	BrName       string `json:"bridge"`
	IsGW         bool   `json:"isGateway"`
	IsDefaultGW  bool   `json:"isDefaultGateway"`
	ForceAddress bool   `json:"forceAddress"`
	IPMasq       bool   `json:"ipMasq"`
	MTU          int    `json:"mtu"`
	HairpinMode  bool   `json:"hairpinMode"`
	OVS          OVS    `json:"ovs"`
}

func loadNetConf(bytes []byte) (*NetConf, error) {
	n := &NetConf{}
	if err := json.Unmarshal(bytes, n); err != nil {
		return nil, fmt.Errorf("failed to load netconf: %v", err)
	}
	return n, nil
}

// vxlanIfName returns formatted vxlan interface name
func vxlanIfName(vtepIP string) string {
	return fmt.Sprintf("vxif%s", strings.Replace(vtepIP, ".", "_", -1))
}

func IPInSlice(ip string, NodeIPs []string) bool {
	for _, i := range NodeIPs {
		if i == ip {
			return true
		}
	}
	return false
}

func setupVTEP(ip string) error {
	// Create interface name for VTEP
	intfName := vxlanIfName(ip)

	// Check if it already exists
	isPresent, vsifName := ovsDriver.IsVtepPresent(ip)
	if !isPresent || (vsifName != intfName) {
		// create VTEP
		err := ovsDriver.CreateVtep(intfName, ip)

		log.Infof("Creating VTEP intf %s for IP %s", intfName, ip)

		if err != nil {
			log.Errorf("Error creating VTEP port %s. Err: %v", intfName, err)
			return err
		}
	}
	return nil
}

func main() {
	raw, e := ioutil.ReadFile("/etc/cni/net.d/linen.conf")
	if e != nil {
		log.Errorf("Read file error: %v\n", e)
		os.Exit(1)
	}

	netConf, err := loadNetConf(raw)
	if err != nil {
		log.Errorf("Load conf error: %v\n", e)
		os.Exit(1)
	}

	if netConf.OVS.IsMaster {
		log.Infof("Daemon running in master")
	} else {
		// Daemon not runs on central node, keep in idle state.
		log.Infof("Daemon running in node")
		for {
			time.Sleep(1 * time.Hour)
		}
	}

	var NodeIPs []string

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// create a ovs bridge
	ovsDriver = ovsdbDriver.NewOvsDriverWithUnix(netConf.OVS.OVSBrName)

	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	for i := 0; i < len(nodes.Items); i++ {
		// fmt.Printf("Initially added nodes %s in to cluster\n", nodes.Items[i].Status.Addresses[0].Address)
		NodeIPs = append(NodeIPs, nodes.Items[i].Status.Addresses[0].Address)
	}

	// monitor
	for {
		nodes, err = clientset.CoreV1().Nodes().List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		// fmt.Printf("There are %d nodes in the cluster\n", len(nodes.Items))

		for i := 0; i < len(nodes.Items); i++ {
			if !IPInSlice(nodes.Items[i].Status.Addresses[0].Address, NodeIPs) {
				if err = setupVTEP(nodes.Items[i].Status.Addresses[0].Address); err != nil {
					log.Errorf("Error creating VTEP port")
				}

				NodeIPs = append(NodeIPs, nodes.Items[i].Status.Addresses[0].Address)
				log.Infof("Added nodes %s in to cluster\n", nodes.Items[i].Status.Addresses[0].Address)
			}
		}
		// fmt.Printf("Current nodes %v\n", NodeIPs)

		time.Sleep(10 * time.Second)
	}
}
