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

package store

import (
	"context"
	"fmt"

	"github.com/mongodb/mongocli/internal/config"
	atlas "go.mongodb.org/atlas/mongodbatlas"
	"go.mongodb.org/ops-manager/opsmngr"
)

//go:generate mockgen -destination=../mocks/mock_clusters.go -package=mocks github.com/mongodb/mongocli/internal/store ClusterLister,AtlasClusterDescriber,OpsManagerClusterDescriber,ClusterCreator,ClusterDeleter,ClusterUpdater,AtlasClusterGetterUpdater,ClusterPauser,ClusterStarter,AtlasClusterQuickStarter,SampleDataAdder,SampleDataStatusDescriber

type ClusterLister interface {
	ProjectClusters(string, *atlas.ListOptions) (interface{}, error)
}

type AtlasClusterDescriber interface {
	AtlasCluster(string, string) (*atlas.Cluster, error)
}

type OpsManagerClusterDescriber interface {
	OpsManagerCluster(string, string) (*opsmngr.Cluster, error)
}

type ClusterCreator interface {
	CreateCluster(*atlas.Cluster) (*atlas.Cluster, error)
}

type ClusterDeleter interface {
	DeleteCluster(string, string) error
}

type ClusterUpdater interface {
	UpdateCluster(string, string, *atlas.Cluster) (*atlas.Cluster, error)
}

type ClusterPauser interface {
	PauseCluster(string, string) (*atlas.Cluster, error)
}

type ClusterStarter interface {
	StartCluster(string, string) (*atlas.Cluster, error)
}

type SampleDataAdder interface {
	AddSampleData(string, string) (*atlas.SampleDatasetJob, error)
}

type SampleDataStatusDescriber interface {
	SampleDataStatus(string, string) (*atlas.SampleDatasetJob, error)
}

type AtlasClusterGetterUpdater interface {
	AtlasClusterDescriber
	ClusterUpdater
}

type AtlasClusterQuickStarter interface {
	SampleDataAdder
	SampleDataStatusDescriber
	CloudProviderRegionsLister
	ClusterLister
	DatabaseUserCreator
	DatabaseUserDescriber
	ProjectIPAccessListCreator
	AtlasClusterDescriber
	ClusterCreator
}

// AddSampleData encapsulate the logic to manage different cloud providers
func (s *Store) AddSampleData(groupID, clusterName string) (*atlas.SampleDatasetJob, error) {
	switch s.service {
	case config.CloudService:
		result, _, err := s.client.(*atlas.Client).Clusters.LoadSampleDataset(context.Background(), groupID, clusterName)
		return result, err
	default:
		return nil, fmt.Errorf("unsupported service: %s", s.service)
	}
}

// SampleData encapsulate the logic to manage different cloud providers
func (s *Store) SampleDataStatus(groupID, id string) (*atlas.SampleDatasetJob, error) {
	switch s.service {
	case config.CloudService:
		result, _, err := s.client.(*atlas.Client).Clusters.GetSampleDatasetStatus(context.Background(), groupID, id)
		return result, err
	default:
		return nil, fmt.Errorf("unsupported service: %s", s.service)
	}
}

// CreateCluster encapsulate the logic to manage different cloud providers
func (s *Store) CreateCluster(cluster *atlas.Cluster) (*atlas.Cluster, error) {
	switch s.service {
	case config.CloudService:
		result, _, err := s.client.(*atlas.Client).Clusters.Create(context.Background(), cluster.GroupID, cluster)
		return result, err
	default:
		return nil, fmt.Errorf("unsupported service: %s", s.service)
	}
}

// UpdateCluster encapsulate the logic to manage different cloud providers
func (s *Store) UpdateCluster(projectID, name string, cluster *atlas.Cluster) (*atlas.Cluster, error) {
	switch s.service {
	case config.CloudService:
		result, _, err := s.client.(*atlas.Client).Clusters.Update(context.Background(), projectID, name, cluster)
		return result, err
	default:
		return nil, fmt.Errorf("unsupported service: %s", s.service)
	}
}

// PauseCluster encapsulate the logic to manage different cloud providers
func (s *Store) PauseCluster(projectID, name string) (*atlas.Cluster, error) {
	paused := true
	cluster := &atlas.Cluster{
		Paused: &paused,
	}
	return s.UpdateCluster(projectID, name, cluster)
}

// StartCluster encapsulate the logic to manage different cloud providers
func (s *Store) StartCluster(projectID, name string) (*atlas.Cluster, error) {
	paused := false
	cluster := &atlas.Cluster{
		Paused: &paused,
	}
	return s.UpdateCluster(projectID, name, cluster)
}

// DeleteCluster encapsulate the logic to manage different cloud providers
func (s *Store) DeleteCluster(projectID, name string) error {
	switch s.service {
	case config.CloudService:
		_, err := s.client.(*atlas.Client).Clusters.Delete(context.Background(), projectID, name)
		return err
	default:
		return fmt.Errorf("unsupported service: %s", s.service)
	}
}

// ProjectClusters encapsulate the logic to manage different cloud providers
func (s *Store) ProjectClusters(projectID string, opts *atlas.ListOptions) (interface{}, error) {
	switch s.service {
	case config.CloudService:
		result, _, err := s.client.(*atlas.Client).Clusters.List(context.Background(), projectID, opts)
		return result, err
	case config.OpsManagerService, config.CloudManagerService:
		result, _, err := s.client.(*opsmngr.Client).Clusters.List(context.Background(), projectID, opts)
		return result, err
	default:
		return nil, fmt.Errorf("unsupported service: %s", s.service)
	}
}

// AtlasCluster encapsulates the logic to manage different cloud providers
func (s *Store) AtlasCluster(projectID, name string) (*atlas.Cluster, error) {
	switch s.service {
	case config.CloudService:
		result, _, err := s.client.(*atlas.Client).Clusters.Get(context.Background(), projectID, name)
		return result, err
	default:
		return nil, fmt.Errorf("unsupported service: %s", s.service)
	}
}

// OpsManagerCluster encapsulates the logic to manage different cloud providers
func (s *Store) OpsManagerCluster(projectID, name string) (*opsmngr.Cluster, error) {
	switch s.service {
	case config.OpsManagerService, config.CloudManagerService:
		result, _, err := s.client.(*opsmngr.Client).Clusters.Get(context.Background(), projectID, name)
		return result, err
	default:
		return nil, fmt.Errorf("unsupported service: %s", s.service)
	}
}

// ListAllProjectClusters encapsulate the logic to manage different cloud providers
func (s *Store) ListAllProjectClusters() (*opsmngr.AllClustersProjects, error) {
	switch s.service {
	case config.OpsManagerService, config.CloudManagerService:
		result, _, err := s.client.(*opsmngr.Client).Clusters.ListAll(context.Background())
		return result, err
	default:
		return nil, fmt.Errorf("unsupported service: %s", s.service)
	}
}
