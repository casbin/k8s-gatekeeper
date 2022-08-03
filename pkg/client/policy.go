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
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8sGateKeeperClient) GetPolicy() ([][]string, error) {
	return k.GetNamedPolicy("p")
}

func (k *K8sGateKeeperClient) GetNamedPolicy(ptype string) ([][]string, error) {
	if k.namespace == "" {
		k.namespace = "default"
	}
	if k.modelName == "" {
		return [][]string{}, fmt.Errorf("modelName should not be empty")
	}
	policyObject, err := k.Clientset.AuthV1().CasbinPolicies(k.namespace).Get(context.TODO(), k.modelName, v1.GetOptions{})
	if err != nil {
		return [][]string{}, err
	}

	policies := stringToPolicies(policyObject.Spec.PolicyItem)

	res := [][]string{}
	for _, p := range policies {
		if len(p) != 0 && p[0] == ptype {
			res = append(res, p[1:])
		}
	}
	return res, nil

}

func (k *K8sGateKeeperClient) AddPolicy(rule []string) (bool, error) {
	return k.AddNamedPolicies("p", [][]string{rule})

}

func (k *K8sGateKeeperClient) AddPolicies(rules [][]string) (bool, error) {
	return k.AddNamedPolicies("p", rules)

}

func (k *K8sGateKeeperClient) AddNamedPolicy(ptype string, rule []string) (bool, error) {
	return k.AddNamedPolicies(ptype, [][]string{rule})
}

func (k *K8sGateKeeperClient) AddNamedPolicies(ptype string, rules [][]string) (bool, error) {
	if k.namespace == "" {
		k.namespace = "default"
	}
	if k.modelName == "" {
		return false, fmt.Errorf("modelName should not be empty")
	}

	policyObject, err := k.Clientset.AuthV1().CasbinPolicies(k.namespace).Get(context.TODO(), k.modelName, v1.GetOptions{})
	if err != nil {
		return false, err
	}

	if policyObject.Spec.PolicyItem != "" && !strings.HasSuffix(policyObject.Spec.PolicyItem, "\n") {
		policyObject.Spec.PolicyItem += "\n"
	}

	for _, rule := range rules {
		policyObject.Spec.PolicyItem += policyToString(ptype, rule...) + "\n"
	}

	_, err = k.Clientset.AuthV1().CasbinPolicies(k.namespace).Update(context.TODO(), policyObject, v1.UpdateOptions{})

	if err != nil {
		return false, err
	}
	return true, nil
}

func (k *K8sGateKeeperClient) RemovePolicy(rule []string) (bool, error) {
	return k.RemoveNamedPolicies("p", [][]string{rule})

}

func (k *K8sGateKeeperClient) RemovePolicies(rules [][]string) (bool, error) {
	return k.RemoveNamedPolicies("p", rules)

}

func (k *K8sGateKeeperClient) RemoveNamedPoliciy(ptype string, rule []string) (bool, error) {
	return k.RemoveNamedPolicies(ptype, [][]string{rule})
}

func (k *K8sGateKeeperClient) RemoveNamedPolicies(ptype string, rules [][]string) (bool, error) {
	if k.namespace == "" {
		k.namespace = "default"
	}
	if k.modelName == "" {
		return false, fmt.Errorf("modelName should not be empty")
	}

	policyObject, err := k.Clientset.AuthV1().CasbinPolicies(k.namespace).Get(context.TODO(), k.modelName, v1.GetOptions{})
	if err != nil {
		return false, err
	}

	policies := stringToPolicies(policyObject.Spec.PolicyItem)

	newPolicies := [][]string{}
	for _, p := range policies {
		remove := false
		for _, r := range rules {
			if len(p) > 0 && policyToString(ptype, r...) == policyToString(p[0], p[1:]...) {
				remove = true
			}
		}

		if !remove {
			newPolicies = append(newPolicies, p)
		}
	}

	policyObject.Spec.PolicyItem = policiesToString(newPolicies)
	_, err = k.Clientset.AuthV1().CasbinPolicies(k.namespace).Update(context.TODO(), policyObject, v1.UpdateOptions{})

	if err != nil {
		return false, err
	}
	return true, nil

}

func (k *K8sGateKeeperClient) UpdatePolicy(oldRule []string, newRule []string) (bool, error) {
	return k.UpdateNamedPolicies("p", [][]string{oldRule}, [][]string{newRule})

}

func (k *K8sGateKeeperClient) UpdatePolicies(oldRules [][]string, newRules [][]string) (bool, error) {
	return k.UpdateNamedPolicies("p", oldRules, newRules)
}

func (k *K8sGateKeeperClient) UpdateNamedPolicy(ptype string, p1 []string, p2 []string) (bool, error) {
	return k.UpdateNamedPolicies(ptype, [][]string{p1}, [][]string{p2})
}

func (k *K8sGateKeeperClient) UpdateNamedPolicies(ptype string, p1 [][]string, p2 [][]string) (bool, error) {
	if len(p1) != len(p2) {
		return false, fmt.Errorf("target policies p1 and new policies p2 should have same length")
	}

	if k.namespace == "" {
		k.namespace = "default"
	}
	if k.modelName == "" {
		return false, fmt.Errorf("modelName should not be empty")
	}

	policyObject, err := k.Clientset.AuthV1().CasbinPolicies(k.namespace).Get(context.TODO(), k.modelName, v1.GetOptions{})
	if err != nil {
		return false, err
	}

	policies := stringToPolicies(policyObject.Spec.PolicyItem)

	newPolicies := [][]string{}

	for _, originalPolicy := range policies {
		match := false
		for i, targetPolicy := range p1 {
			if len(originalPolicy) > 0 && policyToString(ptype, targetPolicy...) == policyToString(originalPolicy[0], originalPolicy[1:]...) {
				match = true
				newPolicies = append(newPolicies, append([]string{ptype}, p2[i]...))
				break
			}
		}

		if !match {
			newPolicies = append(newPolicies, originalPolicy)
		}
	}
	policyObject.Spec.PolicyItem = policiesToString(newPolicies)
	_, err = k.Clientset.AuthV1().CasbinPolicies(k.namespace).Update(context.TODO(), policyObject, v1.UpdateOptions{})

	if err != nil {
		return false, err
	}
	return true, nil
}

func (k *K8sGateKeeperClient) SetPolicy(rule []string) (bool, error) {
	return k.SetNamedPolicies("p", [][]string{rule})

}

func (k *K8sGateKeeperClient) SetPolicies(rules [][]string) (bool, error) {
	return k.SetNamedPolicies("p", rules)

}

func (k *K8sGateKeeperClient) SetNamedPolicy(ptype string, rule []string) (bool, error) {
	return k.SetNamedPolicies(ptype, [][]string{rule})
}

func (k *K8sGateKeeperClient) SetNamedPolicies(ptype string, rules [][]string) (bool, error) {
	if k.namespace == "" {
		k.namespace = "default"
	}
	if k.modelName == "" {
		return false, fmt.Errorf("modelName should not be empty")
	}

	policyObject, err := k.Clientset.AuthV1().CasbinPolicies(k.namespace).Get(context.TODO(), k.modelName, v1.GetOptions{})
	if err != nil {
		return false, err
	}

	policyObject.Spec.PolicyItem = ""

	for _, rule := range rules {
		policyObject.Spec.PolicyItem += policyToString(ptype, rule...) + "\n"
	}

	_, err = k.Clientset.AuthV1().CasbinPolicies(k.namespace).Update(context.TODO(), policyObject, v1.UpdateOptions{})

	if err != nil {
		return false, err
	}
	return true, nil
}
