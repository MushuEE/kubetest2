/*
Copyright 2022 The Kubernetes Authors.

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

package util

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"k8s.io/klog/v2"
	"sigs.k8s.io/kubetest2/pkg/exec"
)

var (
	gkeMinorVersionRegex = regexp.MustCompile(`^v?(\d\.\d+).*$`)
)

// StageGKEBuildMarker stages the build marker to the stage location.
func StageGKEBuildMarker(version, stageLocation, markerPrefix string) error {
	destinationURL := stageLocation + "/" + fileName(version, markerPrefix)
	pushCmd := fmt.Sprintf("gsutil -h 'Content-Type:text/plain' cp - %s", destinationURL)
	cmd := exec.RawCommand(pushCmd)
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}
	cmd.SetStdin(strings.NewReader(version))
	exec.SetOutput(cmd, os.Stdout, os.Stderr)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to upload the latest version number: %s", err)
	}
	return nil
}

func fileName(version, markerPrefix string) string {
	m := gkeMinorVersionRegex.FindStringSubmatch(version)
	if len(m) < 2 {
		klog.Warningf("can't find the minor version of %s, defaulting to %s.txt %v", version, markerPrefix, m)
		return fmt.Sprintf("%s.txt", markerPrefix)
	}

	minor := m[1]
	return fmt.Sprintf("%s-%s.txt", markerPrefix, minor)
}
