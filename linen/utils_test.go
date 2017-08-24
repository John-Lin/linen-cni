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
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestVxlanIfName(t *testing.T) {
	// Test to returns formatted vxlan interface name
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	var reg []string
	for i := 0; i <= 3; i++ {
		reg = append(reg, strconv.Itoa(r1.Intn(256)))
	}

	intfName := vxlanIfName(strings.Join(reg[:], "."))

	checked := "vxif" + strings.Join(reg[:], "_")

	if intfName != checked {
		t.Error("Failed to generate vxlanIfName")
	}
}
