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
	"strconv"
)

func ParseFloat(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ParseFloat requires 1 parameters, currently %d", len(args))
	}

	numString, ok := args[0].(string)
	if ok {
		num, err := strconv.ParseFloat(numString, 64)
		return num, err
	}
	numInt64, ok := args[0].(int64)
	if ok {
		return float64(numInt64), nil
	}
	numInt32, ok := args[0].(int32)
	if ok {
		return float64(numInt32), nil
	}
	numInt, ok := args[0].(int)
	if ok {
		return float64(numInt), nil
	}

	return nil, fmt.Errorf("ParseFloat requires 1st parameter to be string or int")

}

func ParseInt(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ParseFloat requires 1 parameters, currently %d", len(args))
	}
	numString, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("ParseFloat requires 1st parameter to be string")
	}
	num, err := strconv.Atoi(numString)
	return num, err
}

func ToString(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ToString requires 1 parameters, currently %d", len(args))
	}
	v := reflect.ValueOf(args[0])
	if v.Kind() != reflect.String {
		return nil, fmt.Errorf("ToString: args[0] cannot be converted to string")

	}
	return v.String(), nil
}
