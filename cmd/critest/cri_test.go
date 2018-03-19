/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	"github.com/onsi/gomega"

	"github.com/kubernetes-incubator/cri-tools/pkg/framework"

	_ "github.com/kubernetes-incubator/cri-tools/pkg/benchmark"
	_ "github.com/kubernetes-incubator/cri-tools/pkg/validate"
)

var (
	isBenchMark = flag.Bool("benchmark", false, "Run benchmarks instead of validation tests")
	parallel    = flag.Int("parallel", 1, "Run benchmarks in parallel")
)

func init() {
	framework.RegisterFlags()
}

// runTestSuite runs cri validation tests and benchmark tests.
func runTestSuite(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)

	if *isBenchMark {
		flag.Set("ginkgo.focus", "benchmark")
	} else {
		// Skip benchamark measurements for validation tests.
		flag.Set("ginkgo.skipMeasurements", "true")
	}

	reporter := []ginkgo.Reporter{}
	if framework.TestContext.ReportDir != "" {
		if err := os.MkdirAll(framework.TestContext.ReportDir, 0755); err != nil {
			t.Errorf("Failed creating report directory: %v", err)
		}

		reporter = append(reporter, reporters.NewJUnitReporter(path.Join(framework.TestContext.ReportDir, fmt.Sprintf("junit_%v.xml", framework.TestContext.ReportPrefix))))
	}

	ginkgo.RunSpecsWithDefaultAndCustomReporters(t, "CRI validation", reporter)
}

func runParallelTestSuite(t *testing.T) {
	criPath, err := exec.LookPath("critest")
	if err != nil {
		t.Errorf("Failed to lookup path of critest: %v", err)
	}

	args := []string{fmt.Sprintf("-nodes=%d", *parallel), criPath}
	cmd := exec.Command("ginkgo", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Errorf("Failed to run tests in paralllel: %v", err)
	}
}

func TestCRISuite(t *testing.T) {
	if *parallel == 1 {
		runTestSuite(t)
	} else {
		runParallelTestSuite(t)
	}
}
