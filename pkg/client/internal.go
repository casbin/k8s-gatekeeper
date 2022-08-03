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
	"path/filepath"
	"strings"

	"github.com/casbin/k8s-gatekeeper/pkg/generated/clientset/versioned"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func (k *K8sGateKeeperClient) establishInternalClient() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	clientset, err := versioned.NewForConfig(config)
	if err != nil {
		return err
	}
	k.Clientset = clientset
	return nil
}

func (k *K8sGateKeeperClient) establishExternalClient() error {
	home := homedir.HomeDir()
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
	if err != nil {
		return err
	}
	clientset, err := versioned.NewForConfig(config)
	if err != nil {
		return err
	}
	k.Clientset = clientset
	return nil
}
func trim(s string) string {
	return strings.Trim(s, "\r\n ")
}

func policyToString(ptype string, param ...string) string {
	tmp := append([]string{ptype}, param...)
	return strings.Join(tmp, ", ")
}

func stringToPolicies(s string) [][]string {
	res := [][]string{}
	lines := strings.Split(trim(s), "\n")
	for _, s := range lines {
		if s == "" {
			continue
		}
		policy := strings.Split(trim(s), ",")
		for i, t := range policy {
			policy[i] = trim(t)
		}
		res = append(res, policy)
	}
	return res
}

func policiesToString(policies [][]string) string {
	var buffer strings.Builder
	for _, p := range policies {
		if len(p) > 0 {
			buffer.WriteString(policyToString(p[0], p[1:]...))
			buffer.WriteString("\n")
		}
	}
	return buffer.String()
}
