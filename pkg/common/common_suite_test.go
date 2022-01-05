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

package common

import (
	"context"
	stdlog "log"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	policyv1 "github.com/stolostron/cert-policy-controller/api/v1"
)

var cfg *rest.Config

func TestMain(m *testing.M) {
	t := &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "..", "deploy", "crds")},
	}

	if err := policyv1.AddToScheme(scheme.Scheme); err != nil {
		stdlog.Fatal(err)
	}

	var err error
	if cfg, err = t.Start(); err != nil {
		stdlog.Fatal(err)
	}

	code := m.Run()

	if err := t.Stop(); err != nil {
		stdlog.Fatal(err)
	}

	os.Exit(code)
}

// StartTestManager adds recFn.
func StartTestManager(mgr manager.Manager, g *gomega.GomegaWithT) (context.CancelFunc, *sync.WaitGroup) {
	ctx, stop := context.WithCancel(context.TODO())
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		g.Expect(mgr.Start(ctx)).NotTo(gomega.HaveOccurred())
		wg.Done()
	}()

	return stop, wg
}
