// Copyright (C) 2020 - present MongoDB, Inc.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the Server Side Public License, version 1,
// as published by MongoDB, Inc.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// Server Side Public License for more details.
//
// You should have received a copy of the Server Side Public License
// along with this program. If not, see
// http://www.mongodb.com/licensing/server-side-public-license
//
// As a special exception, the copyright holders give permission to link the
// code of portions of this program with the OpenSSL library under certain
// conditions as described in each individual source file and distribute
// linked combinations including the program with the OpenSSL library. You
// must comply with the Server Side Public License in all respects for
// all of the code used other than as permitted herein. If you modify file(s)
// with this exception, you may extend this exception to your version of the
// file(s), but you are not obligated to do so. If you do not wish to do so,
// delete this exception statement from your version. If you delete this
// exception statement from all source files in the program, then also delete
// it in the license file.

package cli

import (
	"github.com/10gen/mcli/internal/flags"
	"github.com/10gen/mcli/internal/json"
	"github.com/10gen/mcli/internal/store"
	"github.com/10gen/mcli/internal/usage"
	"github.com/spf13/cobra"
)

type iamProjectsListOpts struct {
	orgID string
	store store.ProjectLister
}

func (opts *iamProjectsListOpts) init() error {
	s, err := store.New()

	if err != nil {
		return err
	}

	opts.store = s
	return nil
}

func (opts *iamProjectsListOpts) Run() error {
	var projects interface{}
	var err error
	if opts.orgID == "" {
		projects, err = opts.store.GetAllProjects()
	} else {
		projects, err = opts.store.GetOrgProjects(opts.orgID)
	}
	if err != nil {
		return err
	}
	return json.PrettyPrint(projects)
}

// mcli iam project(s) list [--orgId orgId]
func IAMProjectsListBuilder() *cobra.Command {
	opts := new(iamProjectsListOpts)
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List projects",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.init()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.Run()
		},
	}

	cmd.Flags().StringVar(&opts.orgID, flags.OrgID, "", usage.OrgID)

	return cmd
}
