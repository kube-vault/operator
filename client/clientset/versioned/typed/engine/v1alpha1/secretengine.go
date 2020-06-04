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

package v1alpha1

import (
	"context"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1alpha1 "kubevault.dev/operator/apis/engine/v1alpha1"
	scheme "kubevault.dev/operator/client/clientset/versioned/scheme"
)

// SecretEnginesGetter has a method to return a SecretEngineInterface.
// A group's client should implement this interface.
type SecretEnginesGetter interface {
	SecretEngines(namespace string) SecretEngineInterface
}

// SecretEngineInterface has methods to work with SecretEngine resources.
type SecretEngineInterface interface {
	Create(ctx context.Context, secretEngine *v1alpha1.SecretEngine, opts v1.CreateOptions) (*v1alpha1.SecretEngine, error)
	Update(ctx context.Context, secretEngine *v1alpha1.SecretEngine, opts v1.UpdateOptions) (*v1alpha1.SecretEngine, error)
	UpdateStatus(ctx context.Context, secretEngine *v1alpha1.SecretEngine, opts v1.UpdateOptions) (*v1alpha1.SecretEngine, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.SecretEngine, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.SecretEngineList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.SecretEngine, err error)
	SecretEngineExpansion
}

// secretEngines implements SecretEngineInterface
type secretEngines struct {
	client rest.Interface
	ns     string
}

// newSecretEngines returns a SecretEngines
func newSecretEngines(c *EngineV1alpha1Client, namespace string) *secretEngines {
	return &secretEngines{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the secretEngine, and returns the corresponding secretEngine object, and an error if there is any.
func (c *secretEngines) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.SecretEngine, err error) {
	result = &v1alpha1.SecretEngine{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("secretengines").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of SecretEngines that match those selectors.
func (c *secretEngines) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.SecretEngineList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.SecretEngineList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("secretengines").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested secretEngines.
func (c *secretEngines) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("secretengines").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a secretEngine and creates it.  Returns the server's representation of the secretEngine, and an error, if there is any.
func (c *secretEngines) Create(ctx context.Context, secretEngine *v1alpha1.SecretEngine, opts v1.CreateOptions) (result *v1alpha1.SecretEngine, err error) {
	result = &v1alpha1.SecretEngine{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("secretengines").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(secretEngine).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a secretEngine and updates it. Returns the server's representation of the secretEngine, and an error, if there is any.
func (c *secretEngines) Update(ctx context.Context, secretEngine *v1alpha1.SecretEngine, opts v1.UpdateOptions) (result *v1alpha1.SecretEngine, err error) {
	result = &v1alpha1.SecretEngine{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("secretengines").
		Name(secretEngine.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(secretEngine).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *secretEngines) UpdateStatus(ctx context.Context, secretEngine *v1alpha1.SecretEngine, opts v1.UpdateOptions) (result *v1alpha1.SecretEngine, err error) {
	result = &v1alpha1.SecretEngine{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("secretengines").
		Name(secretEngine.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(secretEngine).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the secretEngine and deletes it. Returns an error if one occurs.
func (c *secretEngines) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("secretengines").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *secretEngines) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("secretengines").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched secretEngine.
func (c *secretEngines) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.SecretEngine, err error) {
	result = &v1alpha1.SecretEngine{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("secretengines").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
