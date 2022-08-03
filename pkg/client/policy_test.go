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
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var policyText = `
p,nginx:1.13.1,allow
p,nginx:1.14.1,deny
`

var policyArray = [][]string{
	{"nginx:1.13.1", "allow"},
	{"nginx:1.14.1", "deny"},
}
var policyArray2 = [][]string{
	{"nginx:1.13.1", "allow"},
	{"nginx:1.14.1", "deny"},
	{"nginx:1.14.2", "deny"},
	{"nginx:1.14.3", "deny"},
}
var policyArray3 = [][]string{
	{"nginx:1.13.1", "allow"},
	{"nginx:1.14.1", "deny"},
	{"nginx:1.15.2", "deny"},
	{"nginx:1.15.3", "deny"},
}

func TestPolicy(t *testing.T) {

	Convey("TestCreateModel", t, func() {
		err := client.Namespace("default").ModelName("allowed-repo").CreateModel(modelText)
		So(err, ShouldBeNil)
	})

	Convey("SetNamedPolicies", t, func() {
		ok, err := client.Namespace("default").ModelName("allowed-repo").SetNamedPolicies("p", policyArray)

		So(err, ShouldBeNil)
		So(ok, ShouldBeTrue)

		list2, err := client.Clientset.AuthV1().CasbinPolicies("default").List(context.TODO(), v1.ListOptions{})
		So(err, ShouldBeNil)
		So(len(list2.Items), ShouldEqual, 1)
		So(list2.Items[0].Name, ShouldEqual, "allowed-repo")
		So(list2.Items[0].Namespace, ShouldEqual, "default")
	})

	Convey("GetNamedPolicy", t, func() {
		res, err := client.Namespace("default").ModelName("allowed-repo").GetNamedPolicy("p")
		So(err, ShouldBeNil)

		So(reflect.DeepEqual(res, policyArray), ShouldBeTrue)
	})

	Convey("AddNamedPolicy", t, func() {
		ok, err := client.Namespace("default").ModelName("allowed-repo").AddNamedPolicies("p", policyArray2[2:4])
		So(err, ShouldBeNil)
		So(ok, ShouldBeTrue)

		res, err := client.Namespace("default").ModelName("allowed-repo").GetNamedPolicy("p")
		So(err, ShouldBeNil)
		So(reflect.DeepEqual(res, policyArray2), ShouldBeTrue)
	})

	Convey("UpdateNamedPolicy", t, func() {
		ok, err := client.Namespace("default").ModelName("allowed-repo").UpdateNamedPolicies("p",policyArray2[2:4],policyArray3[2:4])
		So(err, ShouldBeNil)
		So(ok, ShouldBeTrue)

		res, err := client.Namespace("default").ModelName("allowed-repo").GetNamedPolicy("p")
		So(err, ShouldBeNil)
		So(reflect.DeepEqual(res, policyArray3), ShouldBeTrue)
	})
	Convey("RemoveNamedPolicy", t, func() {
		ok, err := client.Namespace("default").ModelName("allowed-repo").RemoveNamedPolicies("p",policyArray3[2:4])
		So(err, ShouldBeNil)
		So(ok, ShouldBeTrue)

		res, err := client.Namespace("default").ModelName("allowed-repo").GetNamedPolicy("p")
		So(err, ShouldBeNil)
		So(reflect.DeepEqual(res, policyArray), ShouldBeTrue)
	})
	reset()
}
