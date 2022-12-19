// Copyright 2019 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Copyright Contributors to the Open Cluster Management project

// Package admissionpolicy handles admissionpolicy controller logic
package common

import (
	"reflect"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	policiesv1 "github.com/stolostron/cert-policy-controller/api/v1"
)

/*
apiVersion: mcm.ibm.com/v1

	kind: GRCPolicy
	metadata:
		name: GRC-policy
	spec:
		namespaces:
			include: ["default"]
			exclude: ["kube*"]
		remediationAction: enforce # or inform
		conditions:
			ownership: [ReplicaSet, Deployment, DeamonSet, ReplicationController]
*/
var plc = &policiesv1.CertificatePolicy{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "testPolicy",
		Namespace: "default",
	},
	Spec: policiesv1.CertificatePolicySpec{
		RemediationAction: policiesv1.Enforce,
		NamespaceSelector: policiesv1.Target{
			Include: []policiesv1.NonEmptyString{"default"},
			Exclude: []policiesv1.NonEmptyString{"kube*"},
		},
	},
}

var sm = SyncedPolicyMap{
	PolicyMap: make(map[string]*policiesv1.CertificatePolicy),
}

// TestGetObject testing get object in map.
func TestGetObject(t *testing.T) {
	t.Parallel()

	_, found := sm.GetObject("void")
	if found {
		t.Fatalf("expecting found = false, however found = %v", found)
	}

	sm.AddObject("default", plc)

	plc, found := sm.GetObject("default")
	if !found {
		t.Fatalf("expecting found = true, however found = %v", found)
	}

	if !reflect.DeepEqual(plc.Name, "testPolicy") {
		t.Fatalf("expecting plcName = testPolicy, however plcName = %v", plc.Name)
	}
}

func TestAddObject(t *testing.T) {
	t.Parallel()
	sm.AddObject("default", plc)

	plcName, found1 := sm.GetObject("ServiceInstance")
	if found1 {
		t.Fatalf("expecting found = false, however found = %v", found1)
	}

	_, found2 := sm.GetObject("void")
	if found2 {
		t.Fatalf("expecting found = false, however found = %v", found2)
	}

	if !reflect.DeepEqual(plc.Name, "testPolicy") {
		t.Fatalf("expecting plcName = testPolicy, however plcName = %v", plcName)
	}
}

func TestRemoveDataObject(t *testing.T) {
	t.Parallel()
	sm.RemoveObject("void")

	_, found := sm.GetObject("void")
	if found {
		t.Fatalf("expecting found = false, however found = %v", found)
	}
	// remove after adding
	sm.AddObject("default", plc)
	sm.RemoveObject("default")

	_, found = sm.GetObject("default")
	if found {
		t.Fatalf("expecting found = false, however found = %v", found)
	}
}
