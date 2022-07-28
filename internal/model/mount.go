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
	"encoding/json"

	admission "k8s.io/api/admission/v1"
	app "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
)

func MountDeploymentObject(admissionReview *admission.AdmissionReview) error {
	admissionReview.Request.Object.Object = nil
	if len(admissionReview.Request.Object.Raw) != 0 {
		var deploymentObject app.Deployment
		err := json.Unmarshal(admissionReview.Request.Object.Raw, &deploymentObject)
		if err != nil {
			return err
		}
		admissionReview.Request.Object.Object = &deploymentObject
	}

	admissionReview.Request.OldObject.Object = nil
	if len(admissionReview.Request.OldObject.Raw) != 0 {
		var deploymentOldObject app.Deployment
		err := json.Unmarshal(admissionReview.Request.OldObject.Raw, &deploymentOldObject)
		if err != nil {
			return err
		}
		admissionReview.Request.OldObject.Object = &deploymentOldObject
	}
	return nil
}

func MountPodObject(admissionReview *admission.AdmissionReview) error {
	admissionReview.Request.Object.Object = nil
	if len(admissionReview.Request.Object.Raw) != 0 {
		var podObject core.Pod
		err := json.Unmarshal(admissionReview.Request.Object.Raw, &podObject)
		if err != nil {
			return err
		}
		admissionReview.Request.Object.Object = &podObject
	}

	admissionReview.Request.OldObject.Object = nil
	if len(admissionReview.Request.OldObject.Raw) != 0 {
		var podOldObject core.Pod
		err := json.Unmarshal(admissionReview.Request.OldObject.Raw, &podOldObject)
		if err != nil {
			return err
		}
		admissionReview.Request.OldObject.Object = &podOldObject
	}
	return nil
}

func MountServiceObject(admissionReview *admission.AdmissionReview) error {
	admissionReview.Request.Object.Object = nil
	if len(admissionReview.Request.Object.Raw) != 0 {
		var serviceObject core.Service
		err := json.Unmarshal(admissionReview.Request.Object.Raw, &serviceObject)
		if err != nil {
			return err
		}
		admissionReview.Request.Object.Object = &serviceObject
	}

	admissionReview.Request.OldObject.Object = nil
	if len(admissionReview.Request.OldObject.Raw) != 0 {
		var serviceOldObject core.Service
		err := json.Unmarshal(admissionReview.Request.OldObject.Raw, &serviceOldObject)
		if err != nil {
			return err
		}
		admissionReview.Request.OldObject.Object = &serviceOldObject
	}
	return nil
}

func MountIngressObject(admissionReview *admission.AdmissionReview) error {
	admissionReview.Request.Object.Object = nil
	if len(admissionReview.Request.Object.Raw) != 0 {
		var ingressObject networking.Ingress
		err := json.Unmarshal(admissionReview.Request.Object.Raw, &ingressObject)
		if err != nil {
			return err
		}
		admissionReview.Request.Object.Object = &ingressObject
	}

	admissionReview.Request.OldObject.Object = nil
	if len(admissionReview.Request.OldObject.Raw) != 0 {
		var ingressOldObject networking.Ingress
		err := json.Unmarshal(admissionReview.Request.OldObject.Raw, &ingressOldObject)
		if err != nil {
			return err
		}
		admissionReview.Request.OldObject.Object = &ingressOldObject
	}
	return nil
}
