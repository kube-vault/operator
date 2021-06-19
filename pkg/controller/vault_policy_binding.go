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
	"time"

	"kubevault.dev/apimachinery/apis"
	policyapi "kubevault.dev/apimachinery/apis/policy/v1alpha1"
	patchutil "kubevault.dev/apimachinery/client/clientset/versioned/typed/policy/v1alpha1/util"
	pbinding "kubevault.dev/operator/pkg/vault/policybinding"

	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog/v2"
	kmapi "kmodules.xyz/client-go/api/v1"
	core_util "kmodules.xyz/client-go/core/v1"
	"kmodules.xyz/client-go/tools/queue"
)

func (c *VaultController) initVaultPolicyBindingWatcher() {
	c.vplcyBindingInformer = c.extInformerFactory.Policy().V1alpha1().VaultPolicyBindings().Informer()
	c.vplcyBindingQueue = queue.New(policyapi.ResourceKindVaultPolicyBinding, c.MaxNumRequeues, c.NumThreads, c.runVaultPolicyBindingInjector)
	c.vplcyBindingInformer.AddEventHandler(queue.NewReconcilableHandler(c.vplcyBindingQueue.GetQueue()))
	if c.auditor != nil {
		c.vplcyBindingInformer.AddEventHandler(c.auditor.ForGVK(policyapi.SchemeGroupVersion.WithKind(policyapi.ResourceKindVaultPolicyBinding)))
	}
	c.vplcyBindingLister = c.extInformerFactory.Policy().V1alpha1().VaultPolicyBindings().Lister()
}

// runVaultPolicyBindingInjector gets the vault policy binding object indexed by the key from cache
// and initializes, reconciles or garbage collects the vault policy binding as needed.
func (c *VaultController) runVaultPolicyBindingInjector(key string) error {
	obj, exists, err := c.vplcyBindingInformer.GetIndexer().GetByKey(key)
	if err != nil {
		klog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		klog.Warningf("VaultPolicyBinding %s does not exist anymore\n", key)
	} else {
		pb := obj.(*policyapi.VaultPolicyBinding).DeepCopy()
		klog.Infof("Sync/Add/Update for VaultPolicyBinding %s/%s\n", pb.Namespace, pb.Name)

		if pb.DeletionTimestamp != nil {
			if core_util.HasFinalizer(pb.ObjectMeta, apis.Finalizer) {
				// Finalize VaultPolicyBinding
				return c.runPolicyBindingFinalizer(pb)
			} else {
				klog.Infof("Finalizer not found for VaultPolicyBinding %s/%s", pb.Namespace, pb.Name)
			}
		} else {
			if !core_util.HasFinalizer(pb.ObjectMeta, apis.Finalizer) {
				// Add finalizer
				_, _, err := patchutil.PatchVaultPolicyBinding(context.TODO(), c.extClient.PolicyV1alpha1(), pb, func(in *policyapi.VaultPolicyBinding) *policyapi.VaultPolicyBinding {
					in.ObjectMeta = core_util.AddFinalizer(pb.ObjectMeta, apis.Finalizer)
					return in
				}, metav1.PatchOptions{})
				if err != nil {
					return errors.Wrapf(err, "failed to add VaultPolicyBinding finalizer for %s/%s", pb.Namespace, pb.Name)
				}
			}

			pbClient, err := pbinding.NewPolicyBindingClient(c.extClient, c.appCatalogClient, c.kubeClient, pb)
			if err != nil {
				return errors.Wrapf(err, "for VaultPolicyBinding %s/%s", pb.Namespace, pb.Name)
			}

			err = c.reconcilePolicyBinding(pb, pbClient)
			if err != nil {
				return errors.Wrapf(err, "for VaultPolicyBinding %s/%s", pb.Namespace, pb.Name)
			}
		}
	}
	return nil
}

// reconcilePolicyBinding reconciles the vault's policy binding
func (c *VaultController) reconcilePolicyBinding(pb *policyapi.VaultPolicyBinding, pbClient pbinding.PolicyBinding) error {
	// create or update policy
	err := pbClient.Ensure(pb)
	if err != nil {
		_, err2 := patchutil.UpdateVaultPolicyBindingStatus(
			context.TODO(),
			c.extClient.PolicyV1alpha1(),
			pb.ObjectMeta,
			func(status *policyapi.VaultPolicyBindingStatus) *policyapi.VaultPolicyBindingStatus {
				status.Phase = policyapi.PolicyBindingFailed
				status.Conditions = kmapi.SetCondition(status.Conditions, kmapi.Condition{
					Type:               kmapi.ConditionFailed,
					Status:             core.ConditionTrue,
					Reason:             "FailedToEnsurePolicyBinding",
					Message:            err.Error(),
					LastTransitionTime: metav1.NewTime(time.Now()),
				})
				return status
			},
			metav1.UpdateOptions{},
		)
		return utilerrors.NewAggregate([]error{err2, err})
	}

	// update status
	_, err = patchutil.UpdateVaultPolicyBindingStatus(
		context.TODO(),
		c.extClient.PolicyV1alpha1(),
		pb.ObjectMeta,
		func(status *policyapi.VaultPolicyBindingStatus) *policyapi.VaultPolicyBindingStatus {
			status.ObservedGeneration = pb.Generation
			status.Phase = policyapi.PolicyBindingSuccess
			status.Conditions = kmapi.RemoveCondition(status.Conditions, kmapi.ConditionFailed)
			status.Conditions = kmapi.SetCondition(status.Conditions, kmapi.Condition{
				Type:    kmapi.ConditionAvailable,
				Status:  core.ConditionTrue,
				Reason:  "Provisioned",
				Message: "policy binding is ready to use",
			})
			return status
		},
		metav1.UpdateOptions{},
	)
	if err != nil {
		return err
	}

	klog.Infof("Successfully processed VaultPolicyBinding: %s/%s", pb.Namespace, pb.Name)
	return nil
}

func (c *VaultController) runPolicyBindingFinalizer(pb *policyapi.VaultPolicyBinding) error {
	klog.Infof("Processing finalizer for VaultPolicyBinding: %s/%s", pb.Namespace, pb.Name)

	pbClient, err := pbinding.NewPolicyBindingClient(c.extClient, c.appCatalogClient, c.kubeClient, pb)
	// The error could be generated for:
	//   - invalid vaultRef in the spec
	// In this case, the operator should be able to delete the VaultPolicyBinding(ie. remove finalizer).
	// If no error occurred:
	//	- Delete the policy
	if err == nil {
		err = pbClient.Delete(pb)
		if err != nil {
			return errors.Wrap(err, "failed to delete the auth role created for policy binding")
		}
	} else {
		klog.Warningf("Skipping cleanup for VaultPolicyBinding: %s/%s with error: %v", pb.Namespace, pb.Name, err)
	}

	// Remove finalizer
	_, err = patchutil.TryPatchVaultPolicyBinding(context.TODO(), c.extClient.PolicyV1alpha1(), pb, func(in *policyapi.VaultPolicyBinding) *policyapi.VaultPolicyBinding {
		in.ObjectMeta = core_util.RemoveFinalizer(in.ObjectMeta, apis.Finalizer)
		return in
	}, metav1.PatchOptions{})
	if err != nil {
		return errors.Wrapf(err, "failed to remove finalizer for VaultPolicyBinding: %s/%s", pb.Namespace, pb.Name)
	}

	klog.Infof("Removed finalizer for VaultPolicyBinding %s/%s", pb.Namespace, pb.Name)
	return nil
}

// finalizePolicyBinding will delete the policy in vault
func (c *VaultController) finalizePolicyBinding(vPBind *policyapi.VaultPolicyBinding) error {

	out, err := c.extClient.PolicyV1alpha1().VaultPolicyBindings(vPBind.Namespace).Get(context.TODO(), vPBind.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

	pBClient, err := pbinding.NewPolicyBindingClient(c.extClient, c.appCatalogClient, c.kubeClient, out)
	if err != nil {
		return err
	}

	return pBClient.Delete(vPBind)
}
