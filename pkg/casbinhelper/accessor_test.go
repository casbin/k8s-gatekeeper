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
package casbinhelper

import (
	"testing"

	casbin "github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	sadaptor "github.com/qiangmzsx/string-adapter/v2"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAccessorNestedStruct(t *testing.T) {
	type Class2 struct {
		Name string
	}
	type Class1 struct {
		Inner Class2
	}
	obj1 := Class1{
		Inner: Class2{
			Name: "1",
		},
	}
	obj2 := Class1{
		Inner: Class2{
			Name: "2",
		},
	}

	Convey("TestAccessorNestedStruct", t, func() {
		modelTxt := `
[request_definition]
r = obj
	
[policy_definition]
p = obj,eft
	
[policy_effect]
e =some(where (p.eft == allow))
	
[matchers]
m = access(r.obj,"Inner","Name")==p.obj
`

		policyTxt := `
p,1,allow
p,2,deny
`
		model := model.NewModel()
		model.LoadModelFromText(modelTxt)
		adaptor := sadaptor.NewAdapter(policyTxt)
		enforcer1, err := casbin.NewEnforcer(model, adaptor)
		So(err, ShouldBeNil)
		enforcer1.AddFunction("access", Access)

		ok1, err1 := enforcer1.Enforce(&obj1)
		So(err1, ShouldBeNil)
		So(ok1, ShouldBeTrue)

		ok2, err2 := enforcer1.Enforce(&obj2)
		So(err2, ShouldBeNil)
		So(ok2, ShouldBeFalse)

		enforcer2, err := casbin.NewEnforcer(model, adaptor)
		So(err, ShouldBeNil)
		enforcer2.AddFunction("access", AccessWithWildCard)
		ok3, err3 := enforcer2.Enforce(&obj1)
		So(err3, ShouldBeNil)
		So(ok3, ShouldBeTrue)

		ok4, err4 := enforcer2.Enforce(&obj2)
		So(err4, ShouldBeNil)
		So(ok4, ShouldBeFalse)

	})

}

func TestAccessorNestedArray(t *testing.T) {
	type Class2 struct {
		Name string
	}
	type Class1 struct {
		Inner []Class2
	}
	obj1 := Class1{
		Inner: []Class2{{
			Name: "1",
		}},
	}
	obj2 := Class1{
		Inner: []Class2{{
			Name: "2",
		}},
	}
	Convey("TestAccessorNestedStruct", t, func() {
		modelTxt := `
[request_definition]
r = obj
	
[policy_definition]
p = obj,eft
	
[policy_effect]
e =some(where (p.eft == allow))
	
[matchers]
m = access(r.obj,"Inner",0,"Name")==p.obj
`

		policyTxt := `
p,1,allow
p,2,deny
`
		model := model.NewModel()
		model.LoadModelFromText(modelTxt)
		adaptor := sadaptor.NewAdapter(policyTxt)
		e, err := casbin.NewEnforcer(model, adaptor)
		So(err, ShouldBeNil)
		e.AddFunction("access", Access)

		ok1, err1 := e.Enforce(&obj1)
		So(err1, ShouldBeNil)
		So(ok1, ShouldBeTrue)

		ok2, err2 := e.Enforce(&obj2)
		So(err2, ShouldBeNil)
		So(ok2, ShouldBeFalse)

		enforcer2, err := casbin.NewEnforcer(model, adaptor)
		So(err, ShouldBeNil)
		enforcer2.AddFunction("access", AccessWithWildCard)
		ok3, err3 := enforcer2.Enforce(&obj1)
		So(err3, ShouldBeNil)
		So(ok3, ShouldBeTrue)

		ok4, err4 := enforcer2.Enforce(&obj2)
		So(err4, ShouldBeNil)
		So(ok4, ShouldBeFalse)
	})

}

type TestAccessorNestedCallClass2 struct {
	Name string
}

func (c *TestAccessorNestedCallClass2) Func1() string {
	return c.Name
}
func (c TestAccessorNestedCallClass2) Func2() string {
	return c.Name
}

type TestAccessorNestedCallClass1 struct {
	Inner []TestAccessorNestedCallClass2
}

func TestAccessorFunctionCall(t *testing.T) {
	obj1 := TestAccessorNestedCallClass1{
		Inner: []TestAccessorNestedCallClass2{{
			Name: "1",
		}},
	}
	obj2 := TestAccessorNestedCallClass1{
		Inner: []TestAccessorNestedCallClass2{{
			Name: "2",
		}},
	}
	Convey("TestAccessorNestedCall1", t, func() {
		modelTxt := `
[request_definition]
r = obj
	
[policy_definition]
p = obj,eft
	
[policy_effect]
e =some(where (p.eft == allow))
	
[matchers]
m = access(r.obj,"Inner",0,"Func1")==p.obj
`

		policyTxt := `
p,1,allow
p,2,deny
`
		model := model.NewModel()
		model.LoadModelFromText(modelTxt)
		adaptor := sadaptor.NewAdapter(policyTxt)
		e, err := casbin.NewEnforcer(model, adaptor)
		So(err, ShouldBeNil)
		e.AddFunction("access", Access)

		ok1, err1 := e.Enforce(&obj1)
		So(err1, ShouldBeNil)
		So(ok1, ShouldBeTrue)

		ok2, err2 := e.Enforce(&obj2)
		So(err2, ShouldBeNil)
		So(ok2, ShouldBeFalse)

		enforcer2, err := casbin.NewEnforcer(model, adaptor)
		So(err, ShouldBeNil)
		enforcer2.AddFunction("access", AccessWithWildCard)
		ok3, err3 := enforcer2.Enforce(&obj1)
		So(err3, ShouldBeNil)
		So(ok3, ShouldBeTrue)

		ok4, err4 := enforcer2.Enforce(&obj2)
		So(err4, ShouldBeNil)
		So(ok4, ShouldBeFalse)
	})
	Convey("TestAccessorNestedCall2", t, func() {
		modelTxt := `
[request_definition]
r = obj
	
[policy_definition]
p = obj,eft
	
[policy_effect]
e =some(where (p.eft == allow))
	
[matchers]
m = access(r.obj,"Inner",0,"Func2")==p.obj
`

		policyTxt := `
p,1,allow
p,2,deny
`
		model := model.NewModel()
		model.LoadModelFromText(modelTxt)
		adaptor := sadaptor.NewAdapter(policyTxt)
		e, err := casbin.NewEnforcer(model, adaptor)
		So(err, ShouldBeNil)
		e.AddFunction("access", Access)

		ok1, err1 := e.Enforce(&obj1)
		So(err1, ShouldBeNil)
		So(ok1, ShouldBeTrue)

		ok2, err2 := e.Enforce(&obj2)
		So(err2, ShouldBeNil)
		So(ok2, ShouldBeFalse)

		enforcer2, err := casbin.NewEnforcer(model, adaptor)
		So(err, ShouldBeNil)
		enforcer2.AddFunction("access", AccessWithWildCard)
		ok3, err3 := enforcer2.Enforce(&obj1)
		So(err3, ShouldBeNil)
		So(ok3, ShouldBeTrue)

		ok4, err4 := enforcer2.Enforce(&obj2)
		So(err4, ShouldBeNil)
		So(ok4, ShouldBeFalse)
	})

	Convey("TestAccessorNestedCall3", t, func() {

		obj1 := TestAccessorNestedCallClass1{
			Inner: []TestAccessorNestedCallClass2{{
				Name: "1",
			}},
		}
		obj2 := TestAccessorNestedCallClass1{
			Inner: []TestAccessorNestedCallClass2{{
				Name: "2",
			}},
		}

		modelTxt := `
[request_definition]
r = obj
	
[policy_definition]
p = obj,eft
	
[policy_effect]
e =some(where (p.eft == allow))
	
[matchers]
m = access(r.obj,"Inner",0,"Func3")==p.obj
`

		policyTxt := `
p,1,allow
p,2,deny
`
		model := model.NewModel()
		model.LoadModelFromText(modelTxt)
		adaptor := sadaptor.NewAdapter(policyTxt)
		e, err := casbin.NewEnforcer(model, adaptor)
		So(err, ShouldBeNil)
		e.AddFunction("access", Access)

		_, err1 := e.Enforce(&obj1)
		So(err1, ShouldNotBeNil)

		_, err2 := e.Enforce(&obj2)
		So(err2, ShouldNotBeNil)

		enforcer2, err := casbin.NewEnforcer(model, adaptor)
		So(err, ShouldBeNil)
		enforcer2.AddFunction("access", AccessWithWildCard)
		_, err3 := enforcer2.Enforce(&obj1)
		So(err3, ShouldNotBeNil)

		_, err4 := enforcer2.Enforce(&obj2)
		So(err4, ShouldNotBeNil)

	})

}

type TestAccessorMapClass2 struct {
	Name string
}

func (t *TestAccessorMapClass2) Func1() string {
	return t.Name
}

type TestAccessorMapClass1 struct {
	Inner map[string]TestAccessorMapClass2
}

func TestAccessorMap(t *testing.T) {

	obj1 := TestAccessorMapClass1{
		Inner: map[string]TestAccessorMapClass2{
			"key": {
				Name: "1",
			}},
	}
	obj2 := TestAccessorMapClass1{
		Inner: map[string]TestAccessorMapClass2{
			"key": {
				Name: "2",
			}},
	}
	Convey("TestAccessorNestedStruct", t, func() {

		modelTxt := `
[request_definition]
r = obj
	
[policy_definition]
p = obj,eft
	
[policy_effect]
e =some(where (p.eft == allow))
	
[matchers]
m = access(r.obj,"Inner","key","Func1")==p.obj
`

		policyTxt := `
p,1,allow
p,2,deny
`
		model := model.NewModel()
		model.LoadModelFromText(modelTxt)
		adaptor := sadaptor.NewAdapter(policyTxt)
		e, err := casbin.NewEnforcer(model, adaptor)
		So(err, ShouldBeNil)
		e.AddFunction("access", Access)

		ok1, err1 := e.Enforce(&obj1)
		So(err1, ShouldBeNil)
		So(ok1, ShouldBeTrue)

		ok2, err2 := e.Enforce(&obj2)
		So(err2, ShouldBeNil)
		So(ok2, ShouldBeFalse)

		enforcer2, err := casbin.NewEnforcer(model, adaptor)
		So(err, ShouldBeNil)
		enforcer2.AddFunction("access", AccessWithWildCard)
		ok3, err3 := enforcer2.Enforce(&obj1)
		So(err3, ShouldBeNil)
		So(ok3, ShouldBeTrue)

		ok4, err4 := enforcer2.Enforce(&obj2)
		So(err4, ShouldBeNil)
		So(ok4, ShouldBeFalse)
	})
	Convey("TestAccessorNestedStruct2", t, func() {

		modelTxt := `
[request_definition]
r = obj
	
[policy_definition]
p = obj,eft
	
[policy_effect]
e =some(where (p.eft == allow))
	
[matchers]
m = access(r.obj,"Inner","key2","Name")==p.obj
`

		policyTxt := `
p,1,allow
p,2,deny
`
		model := model.NewModel()
		model.LoadModelFromText(modelTxt)
		adaptor := sadaptor.NewAdapter(policyTxt)
		e, err := casbin.NewEnforcer(model, adaptor)
		So(err, ShouldBeNil)
		e.AddFunction("access", Access)

		_, err1 := e.Enforce(&obj1)
		So(err1, ShouldNotBeNil)

		_, err2 := e.Enforce(&obj2)
		So(err2, ShouldNotBeNil)

		enforcer2, err := casbin.NewEnforcer(model, adaptor)
		So(err, ShouldBeNil)
		enforcer2.AddFunction("access", AccessWithWildCard)
		_, err3 := enforcer2.Enforce(&obj1)
		So(err3, ShouldNotBeNil)

		_, err4 := enforcer2.Enforce(&obj2)
		So(err4, ShouldNotBeNil)

	})

}

func TestWildcardAccessor(t *testing.T) {
	type Class2 struct {
		Name string
	}
	type Class1 struct {
		Inner []Class2
	}
	type Class0 struct {
		Inner []Class1
	}
	obj1 := Class1{
		Inner: []Class2{
			{Name: "1"},
			{Name: "3"},
		},
	}
	obj2 := Class1{
		Inner: []Class2{
			{Name: "2"},
			{Name: "4"},
		},
	}
	obj3 := Class1{
		Inner: []Class2{
			{Name: "2"},
			{Name: "4"},
		},
	}
	obj4 := Class0{
		Inner: []Class1{obj1, obj3},
	}
	obj5 := Class0{
		Inner: []Class1{obj3, obj2},
	}

	Convey("TestAccessorNestedStruct2", t, func() {

		modelTxt := `
[request_definition]
r = obj
	
[policy_definition]
p = obj,eft
	
[policy_effect]
e =some(where (p.eft == allow))
	
[matchers]
m = contain(access(r.obj,"Inner","*","Inner","*","Name"),p.obj)
`

		policyTxt := `
p,1,allow
p,2,deny
`
		model := model.NewModel()
		model.LoadModelFromText(modelTxt)
		adaptor := sadaptor.NewAdapter(policyTxt)

		enforcer, err := casbin.NewEnforcer(model, adaptor)
		So(err, ShouldBeNil)
		enforcer.AddFunction("access", AccessWithWildCard)
		enforcer.AddFunction("contain", Contain)
		ok1, err1 := enforcer.Enforce(&obj4)
		So(err1, ShouldBeNil)
		So(ok1, ShouldBeTrue)

		ok2, err2 := enforcer.Enforce(&obj5)
		So(err2, ShouldBeNil)
		So(ok2, ShouldBeFalse)
	})

	Convey("TestAccessorNestedStruct", t, func() {

		modelTxt := `
[request_definition]
r = obj
	
[policy_definition]
p = obj,eft
	
[policy_effect]
e =some(where (p.eft == allow))
	
[matchers]
m = contain(access(r.obj,"Inner","*","Name"),p.obj)
`

		policyTxt := `
p,1,allow
p,2,deny
`
		model := model.NewModel()
		model.LoadModelFromText(modelTxt)
		adaptor := sadaptor.NewAdapter(policyTxt)

		enforcer, err := casbin.NewEnforcer(model, adaptor)
		So(err, ShouldBeNil)
		enforcer.AddFunction("access", AccessWithWildCard)
		enforcer.AddFunction("contain", Contain)
		ok1, err1 := enforcer.Enforce(&obj1)
		So(err1, ShouldBeNil)
		So(ok1, ShouldBeTrue)

		ok2, err2 := enforcer.Enforce(&obj2)
		So(err2, ShouldBeNil)
		So(ok2, ShouldBeFalse)
	})
}
