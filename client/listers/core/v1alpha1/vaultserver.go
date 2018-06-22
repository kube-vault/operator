/*
Copyright 2018 The Kube Vault Authors.

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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/kubevault/operator/apis/core/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// VaultServerLister helps list VaultServers.
type VaultServerLister interface {
	// List lists all VaultServers in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.VaultServer, err error)
	// VaultServers returns an object that can list and get VaultServers.
	VaultServers(namespace string) VaultServerNamespaceLister
	VaultServerListerExpansion
}

// vaultServerLister implements the VaultServerLister interface.
type vaultServerLister struct {
	indexer cache.Indexer
}

// NewVaultServerLister returns a new VaultServerLister.
func NewVaultServerLister(indexer cache.Indexer) VaultServerLister {
	return &vaultServerLister{indexer: indexer}
}

// List lists all VaultServers in the indexer.
func (s *vaultServerLister) List(selector labels.Selector) (ret []*v1alpha1.VaultServer, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.VaultServer))
	})
	return ret, err
}

// VaultServers returns an object that can list and get VaultServers.
func (s *vaultServerLister) VaultServers(namespace string) VaultServerNamespaceLister {
	return vaultServerNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// VaultServerNamespaceLister helps list and get VaultServers.
type VaultServerNamespaceLister interface {
	// List lists all VaultServers in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.VaultServer, err error)
	// Get retrieves the VaultServer from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.VaultServer, error)
	VaultServerNamespaceListerExpansion
}

// vaultServerNamespaceLister implements the VaultServerNamespaceLister
// interface.
type vaultServerNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all VaultServers in the indexer for a given namespace.
func (s vaultServerNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.VaultServer, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.VaultServer))
	})
	return ret, err
}

// Get retrieves the VaultServer from the indexer for a given namespace and name.
func (s vaultServerNamespaceLister) Get(name string) (*v1alpha1.VaultServer, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("vaultserver"), name)
	}
	return obj.(*v1alpha1.VaultServer), nil
}
