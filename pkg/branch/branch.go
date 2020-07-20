// Copyright Istio Authors. All Rights Reserved.
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

package branch

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/ghodss/yaml"
	"istio.io/release-builder/pkg/model"
)

// Branch will .....
// This assumes the working directory has been setup and sources resolved.
func Branch(manifest model.Manifest, release string, dryrun bool) error {
	if err := writeManifest(manifest, manifest.OutDir()); err != nil {
		return fmt.Errorf("failed to write manifest: %v", err)
	}

	if err := UpdateDependencies(manifest, dryrun); err != nil {
		return fmt.Errorf("failed to create branches: %v", err)
	}

	// if err := CreateBranches(manifest, release, dryrun); err != nil {
	// 	return fmt.Errorf("failed to create branches: %v", err)
	// }

	// if err := UpdateCommonFiles(manifest, release, dryrun); err != nil {
	// 	return fmt.Errorf("failed to create branches: %v", err)
	// }

	return nil
}

// writeManifest will output the manifest to yaml
func writeManifest(manifest model.Manifest, dir string) error {
	yml, err := yaml.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %v", err)
	}
	if err := ioutil.WriteFile(path.Join(dir, "manifest.yaml"), yml, 0640); err != nil {
		return fmt.Errorf("failed to write manifest: %v", err)
	}
	return nil
}
