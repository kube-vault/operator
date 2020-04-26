/*
Copyright The KubeVault Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "kubevault.dev/operator/apis/engine/v1alpha1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeGCPRoles implements GCPRoleInterface
type FakeGCPRoles struct {
	Fake *FakeEngineV1alpha1
	ns   string
}

var gcprolesResource = schema.GroupVersionResource{Group: "engine.kubevault.com", Version: "v1alpha1", Resource: "gcproles"}

var gcprolesKind = schema.GroupVersionKind{Group: "engine.kubevault.com", Version: "v1alpha1", Kind: "GCPRole"}

// Get takes name of the gCPRole, and returns the corresponding gCPRole object, and an error if there is any.
func (c *FakeGCPRoles) Get(name string, options v1.GetOptions) (result *v1alpha1.GCPRole, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(gcprolesResource, c.ns, name), &v1alpha1.GCPRole{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.GCPRole), err
}

// List takes label and field selectors, and returns the list of GCPRoles that match those selectors.
func (c *FakeGCPRoles) List(opts v1.ListOptions) (result *v1alpha1.GCPRoleList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(gcprolesResource, gcprolesKind, c.ns, opts), &v1alpha1.GCPRoleList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.GCPRoleList{ListMeta: obj.(*v1alpha1.GCPRoleList).ListMeta}
	for _, item := range obj.(*v1alpha1.GCPRoleList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested gCPRoles.
func (c *FakeGCPRoles) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(gcprolesResource, c.ns, opts))

}

// Create takes the representation of a gCPRole and creates it.  Returns the server's representation of the gCPRole, and an error, if there is any.
func (c *FakeGCPRoles) Create(gCPRole *v1alpha1.GCPRole) (result *v1alpha1.GCPRole, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(gcprolesResource, c.ns, gCPRole), &v1alpha1.GCPRole{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.GCPRole), err
}

// Update takes the representation of a gCPRole and updates it. Returns the server's representation of the gCPRole, and an error, if there is any.
func (c *FakeGCPRoles) Update(gCPRole *v1alpha1.GCPRole) (result *v1alpha1.GCPRole, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(gcprolesResource, c.ns, gCPRole), &v1alpha1.GCPRole{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.GCPRole), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeGCPRoles) UpdateStatus(gCPRole *v1alpha1.GCPRole) (*v1alpha1.GCPRole, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(gcprolesResource, "status", c.ns, gCPRole), &v1alpha1.GCPRole{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.GCPRole), err
}

// Delete takes name of the gCPRole and deletes it. Returns an error if one occurs.
func (c *FakeGCPRoles) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(gcprolesResource, c.ns, name), &v1alpha1.GCPRole{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeGCPRoles) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(gcprolesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.GCPRoleList{})
	return err
}

// Patch applies the patch and returns the patched gCPRole.
func (c *FakeGCPRoles) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.GCPRole, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(gcprolesResource, c.ns, name, pt, data, subresources...), &v1alpha1.GCPRole{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.GCPRole), err
}