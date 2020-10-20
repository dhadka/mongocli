// Copyright 2020 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package interfaces

import (
	"github.com/mongodb/mongocli/internal/cli"
	"github.com/mongodb/mongocli/internal/config"
	"github.com/mongodb/mongocli/internal/flag"
	"github.com/mongodb/mongocli/internal/store"
	"github.com/mongodb/mongocli/internal/usage"
	"github.com/spf13/cobra"
)

type CreateOpts struct {
	cli.GlobalOpts
	cli.OutputOpts
	store               store.InterfaceEndpointCreator
	privateEndpointID   string
	interfaceEndpointID string
}

func (opts *CreateOpts) initStore() error {
	var err error
	opts.store, err = store.New(config.Default())
	return err
}

var createTemplate = "Interface endpoint '{{.ID}}' created.\n"

func (opts *CreateOpts) Run() error {
	r, err := opts.store.CreateInterfaceEndpoint(opts.ConfigProjectID(), opts.privateEndpointID, opts.interfaceEndpointID)
	if err != nil {
		return err
	}

	return opts.Print(r)
}

// mongocli atlas privateEndpoint(s)|privateendpoint(s) interface(s) create <interfaceEndpointId> [--privateEndpointId privateEndpointID][--projectId projectId]
func CreateBuilder() *cobra.Command {
	opts := &CreateOpts{}
	cmd := &cobra.Command{
		Use:     "create <interfaceEndpointId>",
		Aliases: []string{"add"},
		Short:   createInterfaceEndpoint,
		Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.PreRunE(
				opts.ValidateProjectID,
				opts.initStore,
				opts.InitOutput(cmd.OutOrStdout(), createTemplate),
			)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.interfaceEndpointID = args[0]
			return opts.Run()
		},
	}
	cmd.Flags().StringVar(&opts.privateEndpointID, flag.PrivateEndpointID, "", usage.PrivateEndpointID)

	cmd.Flags().StringVar(&opts.ProjectID, flag.ProjectID, "", usage.ProjectID)
	cmd.Flags().StringVarP(&opts.Output, flag.Output, flag.OutputShort, "", usage.FormatOut)

	_ = cmd.MarkFlagRequired(flag.PrivateEndpointID)

	return cmd
}
