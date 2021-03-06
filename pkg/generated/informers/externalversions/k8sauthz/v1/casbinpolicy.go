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

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	"context"
	time "time"

	k8sauthzv1 "github.com/casbin/k8s-gatekeeper/pkg/apis/k8sauthz/v1"
	versioned "github.com/casbin/k8s-gatekeeper/pkg/generated/clientset/versioned"
	internalinterfaces "github.com/casbin/k8s-gatekeeper/pkg/generated/informers/externalversions/internalinterfaces"
	v1 "github.com/casbin/k8s-gatekeeper/pkg/generated/listers/k8sauthz/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// CasbinPolicyInformer provides access to a shared informer and lister for
// CasbinPolicies.
type CasbinPolicyInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.CasbinPolicyLister
}

type casbinPolicyInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewCasbinPolicyInformer constructs a new informer for CasbinPolicy type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewCasbinPolicyInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredCasbinPolicyInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredCasbinPolicyInformer constructs a new informer for CasbinPolicy type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredCasbinPolicyInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AuthV1().CasbinPolicies(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AuthV1().CasbinPolicies(namespace).Watch(context.TODO(), options)
			},
		},
		&k8sauthzv1.CasbinPolicy{},
		resyncPeriod,
		indexers,
	)
}

func (f *casbinPolicyInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredCasbinPolicyInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *casbinPolicyInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&k8sauthzv1.CasbinPolicy{}, f.defaultInformer)
}

func (f *casbinPolicyInformer) Lister() v1.CasbinPolicyLister {
	return v1.NewCasbinPolicyLister(f.Informer().GetIndexer())
}
