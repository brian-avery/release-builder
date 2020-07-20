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

	"istio.io/pkg/log"
	"istio.io/release-builder/pkg/model"
	"istio.io/release-builder/pkg/util"
)

// UpdateCommonFiles goes to each repo and runs the command to update the common files.
// A prereq for this is that the common-files relase branch has been updated with a
// new UPDATE_BRANCH and image in it's files.
func UpdateCommonFiles(manifest model.Manifest, release string, dryrun bool) error {
	for repo, dep := range manifest.Dependencies.Get() {
		if dep == nil {
			// Missing a dependency is not always a failure; many are optional dependencies just for
			// tagging.
			log.Warnf("missing dependency: %v", repo)
			continue
		}
		if repo == "common-files" {
			// Skip repos
			log.Infof("Skipping repo: %v", repo)
			continue
		}
		log.Infof("***Updating the common-files for %s from directory: %s", repo, manifest.RepoDir(repo))
		// cmd := util.VerboseCommand("sed", "-i", "'s/UPDATE_BRANCH ?=.*/UPDATE_BRANCH ?= \"+release+\"/'", "common/Makefile.common.mk")
		// cmd.Dir = manifest.RepoDir(repo)
		// if err := cmd.Run(); err != nil {
		// 	return err
		// }
		env := []string{"UPDATE_BRANCH=" + release}
		if err := util.RunMake(manifest, repo, env, "update-common"); err != nil {
			return fmt.Errorf("failed to make update-common: %v", err)
		}

		if !dryrun {
			// command to create PR with diffs.
			// cmd := util.VerboseCommand("git", "push", "--set-upstream", "origin", release)
			// cmd.Dir = manifest.RepoDir(repo)
			// This is a test section. dryrun is turned on by default
			// if err := cmd.Run(); err != nil {
			// 	return err
			// }
			log.Warnf("WARNING: here without specifying dryrun")
		}
	}
	return nil
}
