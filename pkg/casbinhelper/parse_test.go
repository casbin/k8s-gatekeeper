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
	"math"
	"testing"
)

var eps float64 = 0.00000001

func TestParseFloatWithCorrectInput(t *testing.T) {
	var input = []string{"100", "100.0", "-100"}
	var res = []float64{100, 100, -100}
	for i := 0; i < len(input); i++ {
		r, err := ParseFloat(input[i])
		if err != nil {
			t.Errorf("input %s, get error %v", input[i], err)
			return
		}
		rFloat := r.(float64)
		if math.Abs(rFloat-res[i]) > eps {
			t.Errorf("input %s,get %v,expect %v", input[i], rFloat, res[i])
			return
		}
	}
}
func TestParseFloatWithInvalid(t *testing.T) {
	_, err := ParseFloat()
	if err == nil {
		t.Error("input no parameters, shold have got error")
		return
	}
	_, err = ParseFloat("666", "555")
	if err == nil {
		t.Error("input 2 parameters, shold have got error")
		return
	}
	_, err = ParseFloat("hh666.666")
	if err == nil {
		t.Error("input 666.666cc, shold have got error")
		return
	}
}
