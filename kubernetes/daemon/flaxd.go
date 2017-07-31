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
	"github.com/containernetworking/cni/pkg/types/current"
)

var ovsDriver *ovsdbDriver.OvsDriver

// OVS corresponds to Open vSwitch Bridge plugin options
type OVS struct {
	IsMaster   bool     `json:"isMaster"`
	OVSBrName  string   `json:"ovsBridge"`
	VtepIPs    []string `json:"vtepIPs"`
	Controller string   `json:"controller,omitempty"`
}

type LinenConf struct {
	types.NetConf // You may wish to not nest this type
	RuntimeConfig struct {
		OVS OVS `json:"ovs"`
	} `json:"runtimeConfig"`

	RawPrevResult *map[string]interface{} `json:"prevResult"`
	PrevResult    *current.Result         `json:"-"`
}

// NetConfList describes an ordered list of networks.
type NetConfList struct {
	CNIVersion string `json:"cniVersion,omitempty"`

	Name    string       `json:"name,omitempty"`
	Plugins []*LinenConf `json:"plugins,omitempty"`
}

func loadNetConf(bytes []byte) (*NetConfList, error) {
	n := &NetConfList{}
	if err := json.Unmarshal(bytes, n); err != nil {
		return nil, fmt.Errorf("failed to load network configuration: %v", err)
	}
	return n, nil
}

// vxlanIfName returns formatted vxlan interface name
func vxlanIfName(vtepIP string) string {
	return fmt.Sprintf("vxif%s", strings.Replace(vtepIP, ".", "_", -1))
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
	raw, e := ioutil.ReadFile("/etc/cni/net.d/linen.conflist")
	if e != nil {
		log.Errorf("Read file error: %v\n", e)
		os.Exit(1)
	}

	netConf, err := loadNetConf(raw)
	if err != nil {
		log.Errorf("Load conf error: %v\n", e)
		os.Exit(1)
	}

	if netConf.Plugins[1].RuntimeConfig.OVS.IsMaster {
		log.Infof("Daemon running in master")
	} else {
		// Daemon not runs on central node, keep in idle state.
		log.Infof("Daemon running in node")
		for {
			time.Sleep(1 * time.Hour)
		}
	}

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
	ovsDriver = ovsdbDriver.NewOvsDriverWithUnix(netConf.Plugins[1].RuntimeConfig.OVS.OVSBrName)

	// monitoring
	for {
		nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		for i := 0; i < len(nodes.Items); i++ {
			log.Infof("Added nodes %s in to cluster\n", nodes.Items[i].Status.Addresses[0].Address)

			if err = setupVTEP(nodes.Items[i].Status.Addresses[0].Address); err != nil {
				log.Errorf("Error creating VTEP port")
			}

		}

		time.Sleep(10 * time.Second)
	}
}
