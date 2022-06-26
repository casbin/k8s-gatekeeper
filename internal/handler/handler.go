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

package handler

import (
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	admission "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

const (
	admissionApiVersion = "admission.k8s.io/v1"
	admissionKind       = "AdmissionReview"
)

var decoder runtime.Decoder

func init() {
	decoder = serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer()
}

//Main Handler
func Handler(c *gin.Context) {

	data, _ := ioutil.ReadAll(c.Request.Body)
	var admissionReview admission.AdmissionReview

	decoder.Decode(data, nil, &admissionReview)

	//for development only. 
	//Todo:remove this block of code
	if admissionReview.Request.Namespace != "default" {
		approveResponse(c, string(admissionReview.Request.UID))
		return
	}
	//fmt.Println(string(data))
	//fmt.Printf("%s\n", admissionReview.Request.Resource.String())
	//currently we are going to handle these resources:
	uid := admissionReview.Request.UID

	fmt.Println("approved")
	approveResponse(c, string(uid))

}

func rejectResponse(c *gin.Context, uid string, rejectReason string) {
	c.JSON(200, gin.H{
		"apiVersion": admissionApiVersion,
		"kind":       admissionKind,
		"response": map[string]interface{}{
			"uid":     uid,
			"allowed": false,
			"status": map[string]interface{}{
				"code":    403,
				"message": rejectReason,
			},
		},
	})
}

func approveResponse(c *gin.Context, uid string) {
	c.JSON(200, gin.H{
		"apiVersion": admissionApiVersion,
		"kind":       admissionKind,
		"response": map[string]interface{}{
			"uid":     uid,
			"allowed": true,
		},
	})
}
