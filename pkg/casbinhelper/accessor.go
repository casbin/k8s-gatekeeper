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
	"fmt"
	"reflect"
)

func Access(args ...interface{}) (interface{}, error) {
	vCurrent := reflect.ValueOf(args[0])
	for pos, field := range args {
		if pos == 0 {
			continue
		}
		if vCurrent.Kind() == reflect.Pointer {
			vCurrent = vCurrent.Elem()
		}

		if vCurrent.Kind() == reflect.Array || vCurrent.Kind() == reflect.Slice {
			indexFloat, ok := field.(float64)
			if !ok {
				return nil, fmt.Errorf("index for a slice should be a integer")
			}
			index := int(indexFloat)
			vCurrent = vCurrent.Index(index)
			continue
		}
		if vCurrent.Kind() == reflect.Struct {
			attr, ok := field.(string)
			if !ok {
				return nil, fmt.Errorf("string field/method should be applied to a struct")
			}
			//check whether it is a field
			newValueObj := vCurrent.FieldByName(attr)
			if !reflect.ValueOf(newValueObj).IsZero() {
				// is a field
				vCurrent = newValueObj
				continue
			}
			method := vCurrent.MethodByName(attr)
			if !reflect.ValueOf(method).IsZero() {
				// is a method
				if method.Type().NumIn() != 0 || method.Type().NumOut() != 1 {
					return nil, fmt.Errorf("access only support method with no parameters and 1 return value")
				}
				returnValue := method.Call([]reflect.Value{})
				vCurrent = returnValue[0]
				continue
			}

			//maybe a method that requires a pointer receiver?

			if vCurrent.CanAddr() {
				method = vCurrent.Addr().MethodByName(attr)
				if !reflect.ValueOf(method).IsZero() {
					// is a method
					if method.Type().NumIn() != 0 || method.Type().NumOut() != 1 {
						return nil, fmt.Errorf("access only support method with no parameters and 1 return value")
					}
					returnValue := method.Call([]reflect.Value{})
					vCurrent = returnValue[0]
					continue
				}
			}
			//maybe new a new object?

			return nil, fmt.Errorf("no attribute/method %s found", attr)
		}

		if vCurrent.Kind() == reflect.Map {
			vField := reflect.ValueOf(field)
			if !vField.CanConvert(vCurrent.Type().Key()) {
				return nil, fmt.Errorf("key %v cannot be converted to %s", field, vCurrent.Type().Key().String())
			}

			vValue := vCurrent.MapIndex(vField.Convert(vCurrent.Type().Key()))
			if reflect.ValueOf(vValue).IsZero() {
				return nil, fmt.Errorf("key %v not found", field)
			}
			vNewObjPtr := reflect.New(vValue.Type())
			vNewObjPtr.Elem().Set(vValue)
			vCurrent = vNewObjPtr.Elem()

			continue
		}

		return nil, fmt.Errorf("unable to process %s", vCurrent.Type().String())

	}
	if vCurrent.Kind() == reflect.Pointer && !vCurrent.IsNil() {
		vCurrent = vCurrent.Elem()
	}

	return vCurrent.Interface(), nil
}

func AccessWithWildCard(args ...interface{}) (interface{}, error) {
	vCurrent := reflect.ValueOf(args[0])
	return accessWithWildCard(vCurrent, args[1:]...)
}

func accessWithWildCard(vCurrent reflect.Value, args ...interface{}) (interface{}, error) {
	if len(args) == 0 {
		return vCurrent.Interface(), nil
	}

	field := args[0]
	if vCurrent.Kind() == reflect.Pointer {
		vCurrent = vCurrent.Elem()
	}
	if vCurrent.Kind() == reflect.Array || vCurrent.Kind() == reflect.Slice {
		wildCard, ok := field.(string)
		if ok && wildCard == "*" {
			arrayLen := vCurrent.Len()
			resArray := []interface{}{}
			for i := 0; i < arrayLen; i++ {
				newTarget := vCurrent.Index(i)
				res, err := accessWithWildCard(newTarget, args[1:]...)
				if err != nil {
					return nil, err
				}
				if tmp, ok := res.([]interface{}); ok {
					resArray = append(resArray, tmp...)
				} else {
					resArray = append(resArray, res)
				}
			}
			return resArray, nil
		}
		indexFloat, ok := field.(float64)
		if !ok {
			return nil, fmt.Errorf("index for a slice should be a integer")
		}
		index := int(indexFloat)
		vCurrent = vCurrent.Index(index)
		return accessWithWildCard(vCurrent, args[1:]...)
	}
	if vCurrent.Kind() == reflect.Struct {
		fmt.Println(vCurrent.Type().String())
		attr, ok := field.(string)
		if !ok {
			return nil, fmt.Errorf("string field/method should be applied to a struct")
		}
		//check whether it is a field
		newValueObj := vCurrent.FieldByName(attr)
		if !reflect.ValueOf(newValueObj).IsZero() {
			// is a field
			vCurrent = newValueObj
			return accessWithWildCard(vCurrent, args[1:]...)
		}
		method := vCurrent.MethodByName(attr)
		if !reflect.ValueOf(method).IsZero() {
			// is a method
			if method.Type().NumIn() != 0 || method.Type().NumOut() != 1 {
				return nil, fmt.Errorf("access only support method with no parameters and 1 return value")
			}
			returnValue := method.Call([]reflect.Value{})
			vCurrent = returnValue[0]
			return accessWithWildCard(vCurrent, args[1:]...)
		}

		//maybe a method that requires a pointer receiver?
		if vCurrent.CanAddr() {
			method = vCurrent.Addr().MethodByName(attr)
			if !reflect.ValueOf(method).IsZero() {
				// is a method
				if method.Type().NumIn() != 0 || method.Type().NumOut() != 1 {
					return nil, fmt.Errorf("access only support method with no parameters and 1 return value")
				}
				returnValue := method.Call([]reflect.Value{})
				vCurrent = returnValue[0]
				return accessWithWildCard(vCurrent, args[1:]...)
			}
		}
		//maybe new a new object?
		return nil, fmt.Errorf("no attribute/method %s found", attr)
	}

	if vCurrent.Kind() == reflect.Map {
		vKey := reflect.ValueOf(field)
		if !vKey.CanConvert(vCurrent.Type().Key()) {
			return nil, fmt.Errorf("key %v cannot be converted to %s", field, vCurrent.Type().Key().String())
		}

		vValue := vCurrent.MapIndex(vKey.Convert(vCurrent.Type().Key()))
		if reflect.ValueOf(vValue).IsZero() {
			return nil, fmt.Errorf("key %v not found", field)
		}
		//copy a new object in case of the need to call methods who use a pointer receiver
		vNewObjPtr := reflect.New(vValue.Type())
		vNewObjPtr.Elem().Set(vValue)
		vCurrent = vNewObjPtr.Elem()
		return accessWithWildCard(vCurrent, args[1:]...)
	}
	return nil, fmt.Errorf("unable to process %s", vCurrent.Type().String())

}

func Contain(args ...interface{}) (interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("Contain requires more than 1 parameters, currently %d", len(args))
	}
	array := args[0 : len(args)-1]
	target := args[len(args)-1]
	for _, obj := range array {
		if reflect.DeepEqual(obj, target) {
			return true, nil
		}
	}
	return false, nil
}
