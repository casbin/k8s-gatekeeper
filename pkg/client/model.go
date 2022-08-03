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

	k8sauthzv1 "github.com/casbin/k8s-gatekeeper/pkg/apis/k8sauthz/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8sGateKeeperClient) CreateModel(modelText string) error {
	if k.namespace == "" {
		k.namespace = "default"
	}
	if k.modelName == "" {
		return fmt.Errorf("modelName should not be empty")
	}
	casbinModel := k8sauthzv1.CasbinModel{}
	casbinModel.ObjectMeta.Namespace = k.namespace
	casbinModel.ObjectMeta.Name = k.modelName
	casbinModel.Spec.Enabled = true
	casbinModel.Spec.ModelText = modelText
	_, err := k.Clientset.AuthV1().CasbinModels(k.namespace).Create(context.TODO(), &casbinModel, v1.CreateOptions{})
	if err != nil {
		return err
	}

	casbinPolicy := k8sauthzv1.CasbinPolicy{}
	casbinPolicy.ObjectMeta.Namespace = k.namespace
	casbinPolicy.ObjectMeta.Name = k.modelName
	_, err = k.Clientset.AuthV1().CasbinPolicies(k.namespace).Create(context.TODO(), &casbinPolicy, v1.CreateOptions{})
	return err

}

func (k *K8sGateKeeperClient) GetModel() (k8sauthzv1.CasbinModelSpec, error) {
	casbinModel, err := k.Clientset.AuthV1().CasbinModels(k.namespace).Get(context.TODO(), k.modelName, v1.GetOptions{})
	return casbinModel.Spec, err
}

func (k *K8sGateKeeperClient) DeleteModel() error {
	if k.namespace == "" {
		k.namespace = "default"
	}
	if k.modelName == "" {
		return fmt.Errorf("modelName should not be empty")
	}
	err := k.Clientset.AuthV1().CasbinModels(k.namespace).Delete(context.TODO(), k.modelName, v1.DeleteOptions{})
	if err != nil {
		return err
	}
	err = k.Clientset.AuthV1().CasbinPolicies(k.namespace).Delete(context.TODO(), k.modelName, v1.DeleteOptions{})
	return err
}

func (k *K8sGateKeeperClient) UpdateModel(spec k8sauthzv1.CasbinModelSpec) error {
	if k.namespace == "" {
		k.namespace = "default"
	}
	if k.modelName == "" {
		return fmt.Errorf("modelName should not be empty")
	}
	casbinModel, err := k.Clientset.AuthV1().CasbinModels(k.namespace).Get(context.TODO(), k.modelName, v1.GetOptions{})
	if err != nil {
		return err
	}
	casbinModel.Spec = spec
	_, err = k.Clientset.AuthV1().CasbinModels(k.namespace).Update(context.TODO(), casbinModel, v1.UpdateOptions{})
	return err
}
