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

package util

import (
	"bytes"

	"istio.io/pkg/log"
	"istio.io/release-builder/pkg/model"
)

// CreateCommit will look for changes. If changes exist, it will create
// a branch and push a commit with the specified commit text
func CreateCommit(manifest model.Manifest, repo, branch string) error {
	output := bytes.Buffer{}
	cmd := VerboseCommand("git", "status", "--porcelain")
	cmd.Dir = manifest.RepoDir(repo)
	cmd.Stdout = &output
	if err := cmd.Run(); err != nil {
		return err
	}
	log.Infof("**** git status -porcelain returned %v:", output, ".")
	if output.Len() == 0 {
		log.Infof("no changes found to commit")
		return nil
	}

	cmd = VerboseCommand("git", "checkout", "-b", branch)
	cmd.Dir = manifest.RepoDir(repo)
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = VerboseCommand("git", "add", "-A")
	cmd.Dir = manifest.RepoDir(repo)
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = VerboseCommand("git", "commit", "-m", "Run update_dependencies before branching")
	cmd.Dir = manifest.RepoDir(repo)
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = VerboseCommand("git", "push", "--set-upstream", "origin", branch)
	cmd.Dir = manifest.RepoDir(repo)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
