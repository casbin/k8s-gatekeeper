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
	"regexp"
	"strings"
)

func HasPrefix(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("HasPrefix requires 2 parameters, currently %d", len(args))
	}
	str, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("HasPrefix requires 1st parameter to be string")
	}
	prefix, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("HasPrefix requires 2nd parameter to be string")
	}
	return strings.HasPrefix(str, prefix), nil
}

func Split(args ...interface{}) (interface{}, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("split requires 3 parameters, currently %d", len(args))
	}
	sep, ok := args[len(args)-2].(string)
	if !ok {
		return nil, fmt.Errorf("split requires penultimate 2nd parameters to be string")
	}
	posfloat, ok := args[len(args)-1].(float64)
	if !ok {
		return nil, fmt.Errorf("split requires penultimate 1st parameters to be number")

	}
	pos := int(posfloat)

	var res = make([]interface{}, 0)
	for i := 0; i < len(args)-2; i++ {
		str, ok := args[i].(string)
		if !ok {
			return nil, fmt.Errorf("split requires parameter to be string")
		}
		splits := strings.Split(str, sep)
		if len(splits) <= pos {
			return nil, fmt.Errorf("index overflow on string %s", str)
		}
		res = append(res, splits[pos])
	}
	return res, nil
}

func MatchRegex(args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("matchRegex requires 2 parameters, currently %d", len(args))
	}
	regexString, ok := args[len(args)-1].(string)
	if !ok {
		return nil, fmt.Errorf("matchRegex requires penultimate 1st parameters to be string")
	}
	regex, err := regexp.Compile(regexString)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(args)-1; i++ {
		str, ok := args[i].(string)
		if !ok {
			return nil, fmt.Errorf("matchRegex requires parameter to be string")
		}
		if !regex.MatchString(str) {
			return false, fmt.Errorf("string %s doesn't match regex %s", str, regex)
		}
	}
	return true, nil

}
