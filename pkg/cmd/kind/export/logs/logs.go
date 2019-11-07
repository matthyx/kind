/*
Copyright 2018 The Kubernetes Authors.

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

// Package logs implements the `logs` command
package logs

import (
	"fmt"

	"github.com/spf13/cobra"

	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cmd"
	"sigs.k8s.io/kind/pkg/fs"
	"sigs.k8s.io/kind/pkg/log"
)

type flagpole struct {
	Name string
}

// NewCommand returns a new cobra.Command for getting the cluster logs
func NewCommand(logger log.Logger, streams cmd.IOStreams) *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args: cobra.MaximumNArgs(1),
		// TODO(bentheelder): more detailed usage
		Use:   "logs [output-dir]",
		Short: "exports logs to a tempdir or [output-dir] if specified",
		Long:  "exports logs to a tempdir or [output-dir] if specified",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(logger, flags, args)
		},
	}
	cmd.Flags().StringVar(&flags.Name, "name", cluster.DefaultName, "the cluster context name")
	return cmd
}

func runE(logger log.Logger, flags *flagpole, args []string) error {
	provider := cluster.NewProvider(
		cluster.ProviderWithLogger(logger),
	)

	// Check if the cluster has any running nodes
	nodes, err := provider.ListNodes(flags.Name)
	if err != nil {
		return err
	}
	if len(nodes) == 0 {
		return fmt.Errorf("unknown cluster %q", flags.Name)
	}

	// get the optional directory argument, or create a tempdir
	var dir string
	if len(args) == 0 {
		t, err := fs.TempDir("", "")
		if err != nil {
			return err
		}
		dir = t
	} else {
		dir = args[0]
	}

	// collect the logs
	if err := provider.CollectLogs(flags.Name, dir); err != nil {
		return err
	}

	logger.V(0).Info("Exported logs to: " + dir)
	return nil
}