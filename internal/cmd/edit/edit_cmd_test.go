/*
Copyright (c) 2025 Red Hat Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the
License. You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific
language governing permissions and limitations under the License.
*/

package edit

import (
	ffv1 "github.com/innabox/fulfillment-common/api/fulfillment/v1"
	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/ginkgo/v2/dsl/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Edit command", func() {
	DescribeTable("isClusterChanging",
		func(cluster *ffv1.Cluster, expected bool) {
			runner := &runnerContext{}
			result := runner.isClusterChanging(cluster)
			Expect(result).To(Equal(expected))
		},
		Entry("no change - spec matches status",
			&ffv1.Cluster{
				Spec: &ffv1.ClusterSpec{
					NodeSets: map[string]*ffv1.ClusterNodeSet{
						"worker": {Size: 3},
					},
				},
				Status: &ffv1.ClusterStatus{
					NodeSets: map[string]*ffv1.ClusterNodeSet{
						"worker": {Size: 3},
					},
				},
			},
			false,
		),
		Entry("change in progress - size mismatch",
			&ffv1.Cluster{
				Spec: &ffv1.ClusterSpec{
					NodeSets: map[string]*ffv1.ClusterNodeSet{
						"worker": {Size: 5},
					},
				},
				Status: &ffv1.ClusterStatus{
					NodeSets: map[string]*ffv1.ClusterNodeSet{
						"worker": {Size: 3},
					},
				},
			},
			true,
		),
		Entry("change in progress - node set missing in status",
			&ffv1.Cluster{
				Spec: &ffv1.ClusterSpec{
					NodeSets: map[string]*ffv1.ClusterNodeSet{
						"worker": {Size: 3},
					},
				},
				Status: &ffv1.ClusterStatus{
					NodeSets: map[string]*ffv1.ClusterNodeSet{},
				},
			},
			true,
		),
		Entry("no change - multiple node sets all match",
			&ffv1.Cluster{
				Spec: &ffv1.ClusterSpec{
					NodeSets: map[string]*ffv1.ClusterNodeSet{
						"worker":  {Size: 3},
						"control": {Size: 3},
					},
				},
				Status: &ffv1.ClusterStatus{
					NodeSets: map[string]*ffv1.ClusterNodeSet{
						"worker":  {Size: 3},
						"control": {Size: 3},
					},
				},
			},
			false,
		),
		Entry("change in progress - one of multiple node sets mismatch",
			&ffv1.Cluster{
				Spec: &ffv1.ClusterSpec{
					NodeSets: map[string]*ffv1.ClusterNodeSet{
						"worker":  {Size: 5},
						"control": {Size: 3},
					},
				},
				Status: &ffv1.ClusterStatus{
					NodeSets: map[string]*ffv1.ClusterNodeSet{
						"worker":  {Size: 3},
						"control": {Size: 3},
					},
				},
			},
			true,
		),
		Entry("no change - empty node sets",
			&ffv1.Cluster{
				Spec: &ffv1.ClusterSpec{
					NodeSets: map[string]*ffv1.ClusterNodeSet{},
				},
				Status: &ffv1.ClusterStatus{
					NodeSets: map[string]*ffv1.ClusterNodeSet{},
				},
			},
			false,
		),
	)
})
