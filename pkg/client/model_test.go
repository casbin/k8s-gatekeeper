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
	"testing"

	k8sauthzv1 "github.com/casbin/k8s-gatekeeper/pkg/apis/k8sauthz/v1"
	. "github.com/smartystreets/goconvey/convey"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//These tests should be run when a k8s client is available
var modelText = `
[request_definition]
r =  obj

[policy_definition]
p =  obj,eft

[policy_effect]
e = !some(where (p.eft == deny))

[matchers]
m = r.obj.Request.Namespace == "default" && r.obj.Request.Resource.Resource =="deployments" && \
access(r.obj.Request.Object.Object.Spec.Template.Spec.Containers , 0, "Image") == p.obj`

var modelText2 = `
[request_definition]
r =  obj

[policy_definition]
p =  obj,eft

[policy_effect]
e = !some(where (p.eft == deny))

[matchers]
m = r.obj.Request.Namespace == "default" && r.obj.Request.Resource.Resource =="services" && \
r.obj.Request.Operation != "DELETE" &&\
string(r.obj.Request.Object.Object.Spec.Type)  == p.obj
`

func TestModel(t *testing.T) {
	Convey("TestCreateModel", t, func() {
		err := client.Namespace("default").ModelName("allowed-repo").CreateModel(modelText)
		So(err, ShouldBeNil)

		list, err := client.Clientset.AuthV1().CasbinModels("default").List(context.TODO(), v1.ListOptions{})
		So(err, ShouldBeNil)
		So(len(list.Items), ShouldEqual, 1)
		So(list.Items[0].Name, ShouldEqual, "allowed-repo")
		So(list.Items[0].Namespace, ShouldEqual, "default")
		So(list.Items[0].Spec.Enabled, ShouldBeTrue)
		So(list.Items[0].Spec.ModelText, ShouldEqual, modelText)

		list2, err := client.Clientset.AuthV1().CasbinPolicies("default").List(context.TODO(), v1.ListOptions{})
		So(err, ShouldBeNil)
		So(len(list2.Items), ShouldEqual, 1)
		So(list2.Items[0].Name, ShouldEqual, "allowed-repo")
		So(list2.Items[0].Namespace, ShouldEqual, "default")
	})

	Convey("TestGetModel", t, func() {
		spec, err := client.Namespace("default").ModelName("allowed-repo").GetModel()
		So(err, ShouldBeNil)
		So(spec.Enabled, ShouldBeTrue)
		So(spec.ModelText, ShouldEqual, modelText)
	})

	Convey("TestUpdateModel", t, func() {
		newModel := k8sauthzv1.CasbinModelSpec{
			ModelText: modelText2,
			Enabled:   false,
		}
		err := client.Namespace("default").ModelName("allowed-repo").UpdateModel(newModel)
		So(err, ShouldBeNil)

		spec, err := client.Namespace("default").ModelName("allowed-repo").GetModel()
		So(err, ShouldBeNil)
		So(spec.Enabled, ShouldBeFalse)
		So(spec.ModelText, ShouldEqual, modelText2)
	})

	Convey("TestDeleteModel", t, func() {
		err := client.Namespace("default").ModelName("allowed-repo").DeleteModel()
		So(err, ShouldBeNil)

		list, err := client.Clientset.AuthV1().CasbinModels("default").List(context.TODO(), v1.ListOptions{})
		So(err, ShouldBeNil)
		So(len(list.Items), ShouldEqual, 0)

		list2, err := client.Clientset.AuthV1().CasbinPolicies("default").List(context.TODO(), v1.ListOptions{})
		So(err, ShouldBeNil)
		So(len(list2.Items), ShouldEqual, 0)
	})
	

	reset()

}
