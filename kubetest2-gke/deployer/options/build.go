/*
Copyright 2020 The Kubernetes Authors.

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

package options

import (
	"fmt"
	"os"

	gkeBuild "sigs.k8s.io/kubetest2/kubetest2-gke/deployer/build"
	"sigs.k8s.io/kubetest2/pkg/build"
)

type BuildOptions struct {
	CommonBuildOptions      *build.Options
	UpdateLatestGreenMarker bool   `flag:"~update-latest-green-marker" desc:"When set to true, will update the latest-green-x.y.txt marker on GCS."`
	BuildScript             string `flag:"~build-script" desc:"Only used with the gke_make build. Absolute path to the gke_make build script."`
}

var _ build.Builder = &BuildOptions{}
var _ build.Stager = &BuildOptions{}

func (bo *BuildOptions) Validate() error {
	if bo.CommonBuildOptions.Strategy == string(gkeBuild.GKEMakeStrategy) {
		if bo.BuildScript != "" {
			if _, err := os.Stat(bo.BuildScript); err != nil {
				return fmt.Errorf("failed to validate --build-script, required with --strategy=gke_make: %w", err)
			}
			gkeMake := &gkeBuild.GKEMake{
				RepoRoot:      bo.CommonBuildOptions.RepoRoot,
				BuildScript:   bo.BuildScript,
				VersionSuffix: bo.CommonBuildOptions.VersionSuffix,
				StageLocation: bo.CommonBuildOptions.StageLocation,
				UpdateLatest:  bo.CommonBuildOptions.UpdateLatest,
			}
			bo.CommonBuildOptions.Builder = gkeMake
			bo.CommonBuildOptions.Stager = gkeMake
			return nil
		}
	}
	return bo.CommonBuildOptions.Validate()
}

func (bo *BuildOptions) Build() (string, error) {
	return bo.CommonBuildOptions.Build()
}

func (bo *BuildOptions) Stage(version string) error {
	return bo.CommonBuildOptions.Stage(version)
}
