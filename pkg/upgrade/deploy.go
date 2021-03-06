package upgrade

import (
	"context"

	"github.com/openshift/openshift-azure/pkg/api"
	"github.com/openshift/openshift-azure/pkg/util/managedcluster"
)

func (u *simpleUpgrader) Deploy(ctx context.Context, cs *api.OpenShiftManagedCluster, azuretemplate map[string]interface{}, deployFn api.DeployFn) error {
	err := deployFn(ctx, azuretemplate)
	if err != nil {
		return err
	}

	err = u.InitializeCluster(ctx, cs)
	if err != nil {
		return err
	}
	kc, err := managedcluster.ClientsetFromV1ConfigAndWait(ctx, cs.Config.AdminKubeconfig)
	if err != nil {
		return err
	}

	// ensure that all nodes are ready
	return WaitForNodes(ctx, cs, kc)
}
