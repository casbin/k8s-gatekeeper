// Copyright 2022 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"fmt"
	"os/exec"
	"testing"
)

var client *K8sGateKeeperClient

func TestMain(m *testing.M) {
	var err error
	client, err = NewK8sGateKeeperClient(true)

	if err != nil {
		fmt.Println(err)
		return
	}

	reset()
	m.Run()
}

func reset() {
	exec.Command("kubectl", "delete", "-f", "testdata/auth.casbin.org_casbinmodels.yaml").CombinedOutput()

	exec.Command("kubectl", "delete", "-f", "testdata/auth.casbin.org_casbinpolicies.yaml").CombinedOutput()

	res, err := exec.Command("kubectl", "apply", "-f", "testdata/auth.casbin.org_casbinmodels.yaml").CombinedOutput()
	if err != nil {
		fmt.Println(string(res), err)
		return
	}
	res, err = exec.Command("kubectl", "apply", "-f", "testdata/auth.casbin.org_casbinpolicies.yaml").CombinedOutput()
	if err != nil {
		fmt.Println(string(res), err)
		return
	}
}
