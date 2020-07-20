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

	"github.com/spf13/cobra"
	"istio.io/pkg/log"
	"istio.io/release-builder/pkg"
)

var (
	flags = struct {
		manifest    string
		release     string
		githubtoken string
		dryrun      bool
	}{
		manifest: "example/manifest_branch.yaml",
		dryrun:   true, // Default to dry-run for now
	}
	branchCmd = &cobra.Command{
		Use:          "branch",
		Short:        "creates release branches for Istio",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(0),
		RunE: func(c *cobra.Command, _ []string) error {
			if err := validateFlags(); err != nil {
				return fmt.Errorf("invalid flags: %v", err)
			}

			inManifest, err := pkg.ReadInManifest(flags.manifest)
			if err != nil {
				return fmt.Errorf("failed to unmarshal manifest: %v", err)
			}

			manifest, err := pkg.InputManifestToManifest(inManifest)
			if err != nil {
				return fmt.Errorf("failed to setup manifest: %v", err)
			}

			if err := pkg.SetupWorkDir(manifest.Directory); err != nil {
				return fmt.Errorf("failed to setup work dir: %v", err)
			}

			if err := pkg.Sources(manifest); err != nil {
				return fmt.Errorf("failed to fetch sources: %v", err)
			}
			log.Infof("Fetched all sources and setup working directory at %v", manifest.WorkDir())

			if err := pkg.StandardizeManifest(&manifest); err != nil {
				return fmt.Errorf("failed to standardize manifest: %v", err)
			}

			if err := Branch(manifest, flags.release, flags.dryrun); err != nil {
				return fmt.Errorf("failed to branch: %v", err)
			}

			log.Infof("Branched release at %v to release %s", manifest.OutDir(), flags.release)
			return nil
		},
	}
)

func init() {
	branchCmd.PersistentFlags().StringVar(&flags.manifest, "manifest", flags.manifest,
		"The manifest use to get the repos for the branch cut.")
	branchCmd.PersistentFlags().StringVar(&flags.release, "release", flags.release,
		"The directory with the Istio release binary.")
	branchCmd.PersistentFlags().StringVar(&flags.githubtoken, "githubtoken", flags.githubtoken,
		"The file containing a github token.")
	branchCmd.PersistentFlags().BoolVar(&flags.dryrun, "dryrun", flags.dryrun,
		"Do not run any github commands.")
}

func GetBranchCommand() *cobra.Command {
	return branchCmd
}

func validateFlags() error {
	if flags.release == "" {
		return fmt.Errorf("--release required")
	}
	return nil
}
