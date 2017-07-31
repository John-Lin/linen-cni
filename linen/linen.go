// Copyright (c) 2017 Che Wei, Lin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/containernetworking/cni/pkg/version"
)

const defaultOVSBrName = "br0"

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

// parseConfig parses the supplied configuration (and prevResult) from stdin.
func parseConfig(stdin []byte) (*LinenConf, error) {
	conf := LinenConf{}

	if err := json.Unmarshal(stdin, &conf); err != nil {
		return nil, fmt.Errorf("failed to parse network configuration: %v", err)
	}

	// Parse previous result. Remove this if your plugin is not chained.
	if conf.RawPrevResult != nil {
		resultBytes, err := json.Marshal(conf.RawPrevResult)
		if err != nil {
			return nil, fmt.Errorf("could not serialize prevResult: %v", err)
		}
		res, err := version.NewResult(conf.CNIVersion, resultBytes)
		if err != nil {
			return nil, fmt.Errorf("could not parse prevResult: %v", err)
		}
		conf.RawPrevResult = nil
		conf.PrevResult, err = current.NewResultFromResult(res)
		if err != nil {
			return nil, fmt.Errorf("could not convert result to current version: %v", err)
		}
	}

	// if not give ovs bridge name, set to default br0
	if conf.RuntimeConfig.OVS.OVSBrName == "" {
		conf.RuntimeConfig.OVS.OVSBrName = defaultOVSBrName
	}
	return &conf, nil
}

// cmdAdd is called for ADD requests
func cmdAdd(args *skel.CmdArgs) error {
	netConf, err := parseConfig(args.StdinData)

	if err != nil {
		return err
	}

	if netConf.PrevResult == nil {
		return fmt.Errorf("must be called as chained plugin")
	}

	// Create a Open vSwitch bridge
	_, err = setupOVSBridge(netConf)
	if err != nil {
		return err
	}

	// Add Open vSwitch bridge to linux bridge
	if err = addOVSBridgeToBridge(netConf); err != nil {
		return err
	}

	if len(netConf.RuntimeConfig.OVS.VtepIPs) != 0 {
		// Create VxLAN tunnelings
		if err = setupVTEPs(netConf); err != nil {
			return err
		}
	}

	if netConf.RuntimeConfig.OVS.Controller != "" {
		// Set SDN controller
		if err = setupCtrlerToOVS(netConf); err != nil {
			return err
		}
	}

	// Pass through the result for the next plugin
	return types.PrintResult(netConf.PrevResult, netConf.CNIVersion)
}

// cmdDel is called for DELETE requests
func cmdDel(args *skel.CmdArgs) error {
	netConf, err := parseConfig(args.StdinData)
	if err != nil {
		return fmt.Errorf("failed to parse config: %v", err)
	}
	_ = netConf

	return nil
}

func main() {
	skel.PluginMain(cmdAdd, cmdDel, version.PluginSupports("", "0.1.0", "0.2.0", version.Current()))
}
