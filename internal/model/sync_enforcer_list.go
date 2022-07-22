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
	"fmt"
	"log"
	"sync"
	"time"

	casbin "github.com/casbin/casbin/v2"
	"github.com/casbin/k8s-gatekeeper/pkg/casbinhelper"
	admission "k8s.io/api/admission/v1"
)

type EnforcerWrapper struct {
	Enforcer  *casbin.Enforcer
	ModelName string
}

type SyncEnforcerList struct {
	sync.Mutex
	Enforcers []*EnforcerWrapper
	loader    *ModelLoader
}

var EnforcerList *SyncEnforcerList

func Init() {
	EnforcerList = NewSyncEnforcerList()
}

func NewSyncEnforcerList() *SyncEnforcerList {
	//todo: switch to dynamic configuration
	loader, err := NewModelLoader("default", true)
	if err != nil {
		panic(err)
	}
	res := &SyncEnforcerList{
		Enforcers: make([]*EnforcerWrapper, 0),
		loader:    loader,
	}
	//load all enabled models and rules
	res.loadEnforcer()
	//start auto sync for loaders
	go func() {
		for {
			<-time.Tick(5 * time.Second)
			res.loadEnforcer()
		}
	}()
	return res

}

func (s *SyncEnforcerList) Enforce(admission *admission.AdmissionReview) error {
	s.Lock()
	defer s.Unlock()

	for _, wrapper := range s.Enforcers {
		ok, err := wrapper.Enforcer.Enforce(admission)
		if err != nil {
			return fmt.Errorf("%s rejected the request: %s", wrapper.ModelName, err.Error())
		} else if !ok {
			return fmt.Errorf("%s rejected the request", wrapper.ModelName)
		}
	}
	return nil
}

func (s *SyncEnforcerList) loadEnforcer() {
	s.Lock()
	defer s.Unlock()

	modelAndAdptors, err := s.loader.GetModelAndAdaptors()
	if err != nil {
		log.Printf("error: %s", err.Error())
		return
	}
	s.Enforcers = s.Enforcers[:0]
	for _, tmp := range modelAndAdptors {
		e, err := casbin.NewEnforcer(tmp.Model, tmp.Adaptor)
		if err != nil {
			log.Printf("error: %s", err.Error())
			return
		}
		//todo: setup function list
		e.AddFunction("access", casbinhelper.Access)
		s.Enforcers = append(s.Enforcers, &EnforcerWrapper{Enforcer: e, ModelName: tmp.Name})
	}
	log.Printf("%d enforcers loaded", len(s.Enforcers))
}
