// Copyright 2022 The Casbin Authors. All Rights Reserved.
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
package model

import (
	"context"
	"path/filepath"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/casbin/k8s-gatekeeper/pkg/crdadaptor"
	"github.com/casbin/k8s-gatekeeper/pkg/generated/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type ModelLoader struct {
	namespace        string
	isExternalClient bool
	clientset        *versioned.Clientset
}

type ModelAdaptorPair struct {
	Name    string
	Model   model.Model
	Adaptor persist.Adapter
}

func NewModelLoader(namespace string, isExternalClient bool) (*ModelLoader, error) {
	var res = &ModelLoader{
		namespace:        namespace,
		isExternalClient: isExternalClient,
	}
	var err error
	if isExternalClient {
		err = res.establishExternalClient()
	} else {
		err = res.establishInternalClient()
	}
	return res, err
}

func (m *ModelLoader) GetModelAndAdaptors() ([]ModelAdaptorPair, error) {
	list, err := m.clientset.AuthV1().CasbinModels(m.namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	res := make([]ModelAdaptorPair, 0)
	for _, crdModel := range list.Items {
		casbinModel, err := model.NewModelFromString(crdModel.Spec.ModelText)
		if err != nil {
			return nil, err
		}
		modelName := crdModel.ObjectMeta.Name
		casbinAdaptor, err := crdadaptor.NewK8sAdaptor(m.namespace, modelName, m.isExternalClient)
		if err != nil {
			return nil, err
		}
		res = append(res, ModelAdaptorPair{
			Model:   casbinModel,
			Adaptor: casbinAdaptor,
			Name:    modelName,
		})
	}
	return res, nil
}

func (m *ModelLoader) establishInternalClient() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	clientset, err := versioned.NewForConfig(config)
	if err != nil {
		return err
	}
	m.clientset = clientset
	return nil
}

func (m *ModelLoader) establishExternalClient() error {
	home := homedir.HomeDir()
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
	if err != nil {
		return err
	}
	clientset, err := versioned.NewForConfig(config)
	if err != nil {
		return err
	}
	m.clientset = clientset
	return nil
}
