/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Community License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Community-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	"kubevault.dev/apimachinery/apis"
	api "kubevault.dev/apimachinery/apis/engine/v1alpha1"
	patchutil "kubevault.dev/apimachinery/client/clientset/versioned/typed/engine/v1alpha1/util"
	"kubevault.dev/operator/pkg/vault/role/gcp"

	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog/v2"
	kmapi "kmodules.xyz/client-go/api/v1"
	core_util "kmodules.xyz/client-go/core/v1"
	"kmodules.xyz/client-go/tools/queue"
)

const (
	GCPRolePhaseSuccess    api.GCPRolePhase = "Success"
	GCPRolePhaseProcessing api.GCPRolePhase = "Processing"
)

func (c *VaultController) initGCPRoleWatcher() {
	c.gcpRoleInformer = c.extInformerFactory.Engine().V1alpha1().GCPRoles().Informer()
	c.gcpRoleQueue = queue.New(api.ResourceKindGCPRole, c.MaxNumRequeues, c.NumThreads, c.runGCPRoleInjector)
	c.gcpRoleInformer.AddEventHandler(queue.NewReconcilableHandler(c.gcpRoleQueue.GetQueue()))
	if c.auditor != nil {
		c.gcpRoleInformer.AddEventHandler(c.auditor.ForGVK(api.SchemeGroupVersion.WithKind(api.ResourceKindGCPRole)))
	}
	c.gcpRoleLister = c.extInformerFactory.Engine().V1alpha1().GCPRoles().Lister()
}

func (c *VaultController) runGCPRoleInjector(key string) error {
	obj, exist, err := c.gcpRoleInformer.GetIndexer().GetByKey(key)
	if err != nil {
		klog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exist {
		klog.Warningf("GCPRole %s does not exist anymore", key)

	} else {
		role := obj.(*api.GCPRole).DeepCopy()

		klog.Infof("Sync/Add/Update for GCPRole %s/%s", role.Namespace, role.Name)

		if role.DeletionTimestamp != nil {
			if core_util.HasFinalizer(role.ObjectMeta, apis.Finalizer) {
				return c.runGCPRoleFinalizer(role)
			}
		} else {
			if !core_util.HasFinalizer(role.ObjectMeta, apis.Finalizer) {
				// Add finalizer
				_, _, err := patchutil.PatchGCPRole(context.TODO(), c.extClient.EngineV1alpha1(), role, func(in *api.GCPRole) *api.GCPRole {
					in.ObjectMeta = core_util.AddFinalizer(in.ObjectMeta, apis.Finalizer)
					return in
				}, metav1.PatchOptions{})
				if err != nil {
					return errors.Wrapf(err, "failed to add finalizer for GCPRole: %s/%s", role.Namespace, role.Name)
				}
			}

			// Conditions are empty, when the GCPRole obj is enqueued for the first time.
			// Set status.phase to "Processing".
			if role.Status.Conditions == nil {
				newRole, err := patchutil.UpdateGCPRoleStatus(
					context.TODO(),
					c.extClient.EngineV1alpha1(),
					role.ObjectMeta,
					func(status *api.GCPRoleStatus) *api.GCPRoleStatus {
						status.Phase = GCPRolePhaseProcessing
						return status
					},
					metav1.UpdateOptions{},
				)
				if err != nil {
					return errors.Wrapf(err, "failed to update status for GCPRole: %s/%s", role.Namespace, role.Name)
				}
				role = newRole
			}

			rClient, err := gcp.NewGCPRole(c.kubeClient, c.appCatalogClient, role)
			if err != nil {
				return err
			}

			err = c.reconcileGCPRole(rClient, role)
			if err != nil {
				return errors.Wrapf(err, "failed to reconcile GCPRole: %s/%s", role.Namespace, role.Name)
			}
		}
	}
	return nil
}

// Will do:
//	For vault:
// 	  - configure a GCP role
//    - sync role
func (c *VaultController) reconcileGCPRole(rClient gcp.GCPRoleInterface, role *api.GCPRole) error {
	// create role
	err := rClient.CreateRole()
	if err != nil {
		_, err2 := patchutil.UpdateGCPRoleStatus(
			context.TODO(),
			c.extClient.EngineV1alpha1(),
			role.ObjectMeta, func(status *api.GCPRoleStatus) *api.GCPRoleStatus {
				status.Conditions = kmapi.SetCondition(status.Conditions, kmapi.Condition{
					Type:    kmapi.ConditionFailed,
					Status:  core.ConditionTrue,
					Reason:  "FailedToCreateRole",
					Message: err.Error(),
				})
				return status
			},
			metav1.UpdateOptions{},
		)
		return utilerrors.NewAggregate([]error{err2, errors.Wrap(err, "failed to create role")})
	}

	_, err = patchutil.UpdateGCPRoleStatus(
		context.TODO(),
		c.extClient.EngineV1alpha1(),
		role.ObjectMeta, func(status *api.GCPRoleStatus) *api.GCPRoleStatus {
			status.Phase = GCPRolePhaseSuccess
			status.ObservedGeneration = role.Generation
			status.Conditions = kmapi.RemoveCondition(status.Conditions, kmapi.ConditionFailed)
			status.Conditions = kmapi.SetCondition(status.Conditions, kmapi.Condition{
				Type:    kmapi.ConditionAvailable,
				Status:  core.ConditionTrue,
				Reason:  "Provisioned",
				Message: "role is ready to use",
			})
			return status
		},
		metav1.UpdateOptions{},
	)
	if err != nil {
		return err
	}

	klog.Infof("Successfully processed GCPRole: %s/%s", role.Namespace, role.Name)
	return nil
}

func (c *VaultController) runGCPRoleFinalizer(role *api.GCPRole) error {
	klog.Infof("Processing finalizer for GCPRole %s/%s", role.Namespace, role.Name)

	rClient, err := gcp.NewGCPRole(c.kubeClient, c.appCatalogClient, role)
	// The error could be generated for:
	//   - invalid vaultRef in the spec
	// In this case, the operator should be able to delete the GCPRole (ie. remove finalizer).
	// If no error occurred:
	//	- Delete the gcp role created in vault
	if err == nil {
		statusCode, err := rClient.DeleteRole(role.RoleName())
		// For the following errors, the operator should be
		// able to delete the role obj.
		// 	- 400 - Invalid request, missing or invalid data.
		// 	- 403 - Forbidden, your authentication details are either incorrect, you don't have access to this feature, or - if CORS is enabled - you made a cross-origin request from an origin that is not allowed to make such requests.
		//  - 404 - Invalid path. This can both mean that the path truly doesn't exist or that you don't have permission to view a specific path. We use 404 in some cases to avoid state leakage.
		// return error if it is network error.
		if err != nil && (statusCode/100) != 4 {
			return errors.Wrap(err, "failed to delete gcp role")
		}
	} else {
		klog.Warningf("Skipping cleanup for GCPRole: %s/%s with error: %v", role.Namespace, role.Name, err)
	}

	// remove finalizer
	_, _, err = patchutil.PatchGCPRole(context.TODO(), c.extClient.EngineV1alpha1(), role, func(in *api.GCPRole) *api.GCPRole {
		in.ObjectMeta = core_util.RemoveFinalizer(in.ObjectMeta, apis.Finalizer)
		return in
	}, metav1.PatchOptions{})
	if err != nil {
		return errors.Wrapf(err, "failed to remove finalizer for GCPRole: %s/%s", role.Namespace, role.Name)
	}

	klog.Infof("Removed finalizer for GCPRole: %s/%s", role.Namespace, role.Name)
	return nil
}
