package storagecluster

import (
	"context"
	"testing"

	api "github.com/red-hat-storage/ocs-operator/api/v4/v1"
	"github.com/red-hat-storage/ocs-operator/v4/controllers/defaults"
	cephv1 "github.com/rook/rook/pkg/apis/ceph.rook.io/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestCephFileSystem(t *testing.T) {
	var cases = []struct {
		label                  string
		createRuntimeObjects   bool
		remoteStorageConsumers bool
	}{
		{
			label:                "case 1",
			createRuntimeObjects: false,
		},
	}

	for _, c := range cases {
		var objects []client.Object
		t, reconciler, cr, request := initStorageClusterResourceCreateUpdateTest(t, objects, nil)

		if c.createRuntimeObjects {
			objects = createUpdateRuntimeObjects(t) //nolint:staticcheck //no need to use objects as they update in runtime
		}
		assertCephFileSystem(t, reconciler, cr, request)
	}
}

func assertCephFileSystem(t *testing.T, reconciler StorageClusterReconciler, cr *api.StorageCluster, request reconcile.Request) {
	actualFs := &cephv1.CephFilesystem{
		ObjectMeta: metav1.ObjectMeta{
			Name: "ocsinit-cephfilesystem",
		},
	}
	request.Name = "ocsinit-cephfilesystem"
	err := reconciler.Client.Get(context.TODO(), request.NamespacedName, actualFs)
	assert.NoError(t, err)

	expectedAf, err := reconciler.newCephFilesystemInstances(cr)
	assert.NoError(t, err)

	assert.Equal(t, len(expectedAf[0].OwnerReferences), 1)

	assert.Equal(t, expectedAf[0].ObjectMeta.Name, actualFs.ObjectMeta.Name)
	assert.Equal(t, expectedAf[0].Spec, actualFs.Spec)
}

func TestGetActiveMetadataServers(t *testing.T) {
	var cases = []struct {
		label                         string
		sc                            *api.StorageCluster
		expectedActiveMetadataServers int
	}{
		{
			label:                         "Default case",
			sc:                            &api.StorageCluster{},
			expectedActiveMetadataServers: defaults.CephFSActiveMetadataServers,
		},
		{
			label: "ActiveMetadataServers is set on the StorageCluster CR Spec",
			sc: &api.StorageCluster{
				Spec: api.StorageClusterSpec{
					ManagedResources: api.ManagedResourcesSpec{
						CephFilesystems: api.ManageCephFilesystems{
							ActiveMetadataServers: 2,
						},
					},
				},
			},
			expectedActiveMetadataServers: 2,
		},
	}

	for _, c := range cases {
		t.Logf("Case: %s\n", c.label)
		actualActiveMetadataServers := getActiveMetadataServers(c.sc)
		assert.Equal(t, c.expectedActiveMetadataServers, actualActiveMetadataServers)
	}

}

func TestCephFileSystemDataPools(t *testing.T) {
	mocksc := &api.StorageCluster{}
	mockStorageCluster.DeepCopyInto(mocksc)
	mocksc.Status.FailureDomain = "zone"
	defaultPoolSpec := cephv1.PoolSpec{
		EnableCrushUpdates: true,
		DeviceClass:        mocksc.Status.DefaultCephDeviceClass,
		FailureDomain:      getFailureDomain(mocksc),
		Replicated:         generateCephReplicatedSpec(mocksc, poolTypeData),
	}

	var cases = []struct {
		label             string
		sc                *api.StorageCluster
		expectedDataPools []cephv1.NamedPoolSpec
	}{
		{
			label: "Neither DataPoolSpec nor AdditionalDataPools is set",
			sc:    &api.StorageCluster{},
			expectedDataPools: []cephv1.NamedPoolSpec{
				{
					PoolSpec: defaultPoolSpec,
				},
			},
		},
		{
			label: "DataPoolSpec is set & AdditionalDataPools is not set",
			sc: &api.StorageCluster{
				Spec: api.StorageClusterSpec{
					ManagedResources: api.ManagedResourcesSpec{
						CephFilesystems: api.ManageCephFilesystems{
							DataPoolSpec: &cephv1.PoolSpec{
								DeviceClass: "gold",
								Replicated: cephv1.ReplicatedSpec{
									Size:            2,
									TargetSizeRatio: 0.8,
								},
							},
						},
					},
				},
			},
			expectedDataPools: []cephv1.NamedPoolSpec{
				{
					PoolSpec: cephv1.PoolSpec{
						DeviceClass:        "gold",
						EnableCrushUpdates: true,
						Replicated: cephv1.ReplicatedSpec{
							Size:                     2,
							TargetSizeRatio:          0.8,
							ReplicasPerFailureDomain: defaultPoolSpec.Replicated.ReplicasPerFailureDomain,
						},
						FailureDomain: defaultPoolSpec.FailureDomain,
					},
				},
			},
		},
		{
			label: "DataPoolSpec is not set & One item is set on AdditionalDataPools",
			sc: &api.StorageCluster{
				Spec: api.StorageClusterSpec{
					ManagedResources: api.ManagedResourcesSpec{
						CephFilesystems: api.ManageCephFilesystems{
							AdditionalDataPools: []cephv1.NamedPoolSpec{
								{
									Name: "test-1",
									PoolSpec: cephv1.PoolSpec{
										Replicated: cephv1.ReplicatedSpec{
											Size:            2,
											TargetSizeRatio: 0.3,
										},
									},
								},
							},
						},
					},
				},
			},
			expectedDataPools: []cephv1.NamedPoolSpec{
				{
					PoolSpec: defaultPoolSpec,
				},
				{
					Name: "test-1",
					PoolSpec: cephv1.PoolSpec{
						DeviceClass:        defaultPoolSpec.DeviceClass,
						EnableCrushUpdates: true,
						Replicated: cephv1.ReplicatedSpec{
							Size:                     2,
							TargetSizeRatio:          0.3,
							ReplicasPerFailureDomain: defaultPoolSpec.Replicated.ReplicasPerFailureDomain,
						},
						FailureDomain: defaultPoolSpec.FailureDomain,
					},
				},
			},
		},
		{
			label: "DataPoolSpec is not set & multiple AdditionalDataPools are set",
			sc: &api.StorageCluster{
				Spec: api.StorageClusterSpec{
					ManagedResources: api.ManagedResourcesSpec{
						CephFilesystems: api.ManageCephFilesystems{
							AdditionalDataPools: []cephv1.NamedPoolSpec{
								{
									Name: "test-1",
									PoolSpec: cephv1.PoolSpec{
										DeviceClass: "gold",
									},
								},
								{
									Name: "test-2",
									PoolSpec: cephv1.PoolSpec{
										DeviceClass: "silver",
									},
								},
							},
						},
					},
				},
			},
			expectedDataPools: []cephv1.NamedPoolSpec{
				{
					PoolSpec: defaultPoolSpec,
				},
				{
					Name: "test-1",
					PoolSpec: cephv1.PoolSpec{
						DeviceClass:        "gold",
						EnableCrushUpdates: true,
						Replicated:         defaultPoolSpec.Replicated,
						FailureDomain:      defaultPoolSpec.FailureDomain,
					},
				},
				{
					Name: "test-2",
					PoolSpec: cephv1.PoolSpec{
						DeviceClass:        "silver",
						EnableCrushUpdates: true,
						Replicated:         defaultPoolSpec.Replicated,
						FailureDomain:      defaultPoolSpec.FailureDomain,
					},
				},
			},
		},
		{
			label: "DataPoolSpec is set & multiple AdditionalDataPools are set",
			sc: &api.StorageCluster{
				Spec: api.StorageClusterSpec{
					ManagedResources: api.ManagedResourcesSpec{
						CephFilesystems: api.ManageCephFilesystems{
							DataPoolSpec: &cephv1.PoolSpec{
								DeviceClass: "gold",
								Replicated: cephv1.ReplicatedSpec{
									TargetSizeRatio: 0.1,
								},
							},
							AdditionalDataPools: []cephv1.NamedPoolSpec{
								{
									Name: "test-1",
									PoolSpec: cephv1.PoolSpec{
										DeviceClass: "silver",
										Replicated: cephv1.ReplicatedSpec{
											Size:            2,
											TargetSizeRatio: 0.25,
										},
									},
								},
								{
									Name: "test-2",
									PoolSpec: cephv1.PoolSpec{
										DeviceClass: "bronze",
										Replicated: cephv1.ReplicatedSpec{
											Size:            2,
											TargetSizeRatio: 0.25,
										},
									},
								},
							},
						},
					},
				},
			},
			expectedDataPools: []cephv1.NamedPoolSpec{
				{
					PoolSpec: cephv1.PoolSpec{
						DeviceClass:        "gold",
						EnableCrushUpdates: true,
						Replicated: cephv1.ReplicatedSpec{
							Size:                     defaultPoolSpec.Replicated.Size,
							TargetSizeRatio:          0.1,
							ReplicasPerFailureDomain: defaultPoolSpec.Replicated.ReplicasPerFailureDomain,
						},
						FailureDomain: defaultPoolSpec.FailureDomain,
					},
				},
				{
					Name: "test-1",
					PoolSpec: cephv1.PoolSpec{
						DeviceClass:        "silver",
						EnableCrushUpdates: true,
						Replicated: cephv1.ReplicatedSpec{
							Size:                     2,
							TargetSizeRatio:          0.25,
							ReplicasPerFailureDomain: defaultPoolSpec.Replicated.ReplicasPerFailureDomain,
						},
						FailureDomain: defaultPoolSpec.FailureDomain,
					},
				},
				{
					Name: "test-2",
					PoolSpec: cephv1.PoolSpec{
						DeviceClass:        "bronze",
						EnableCrushUpdates: true,
						Replicated: cephv1.ReplicatedSpec{
							Size:                     2,
							TargetSizeRatio:          0.25,
							ReplicasPerFailureDomain: defaultPoolSpec.Replicated.ReplicasPerFailureDomain,
						},
						FailureDomain: defaultPoolSpec.FailureDomain,
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Logf("Case: %s\n", c.label)
		var objects []client.Object
		t, reconciler, _, _ := initStorageClusterResourceCreateUpdateTest(t, objects, nil)
		c.sc.Status.FailureDomain = "zone"
		filesystem, err := reconciler.newCephFilesystemInstances(c.sc)
		assert.NoError(t, err)
		actualDataPools := filesystem[0].Spec.DataPools
		assert.Equal(t, c.expectedDataPools, actualDataPools)
	}
}
