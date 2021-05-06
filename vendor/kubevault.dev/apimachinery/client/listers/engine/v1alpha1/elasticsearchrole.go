/*
Copyright AppsCode Inc. and Contributors

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
	v1alpha1 "kubevault.dev/apimachinery/apis/engine/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ElasticsearchRoleLister helps list ElasticsearchRoles.
// All objects returned here must be treated as read-only.
type ElasticsearchRoleLister interface {
	// List lists all ElasticsearchRoles in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ElasticsearchRole, err error)
	// ElasticsearchRoles returns an object that can list and get ElasticsearchRoles.
	ElasticsearchRoles(namespace string) ElasticsearchRoleNamespaceLister
	ElasticsearchRoleListerExpansion
}

// elasticsearchRoleLister implements the ElasticsearchRoleLister interface.
type elasticsearchRoleLister struct {
	indexer cache.Indexer
}

// NewElasticsearchRoleLister returns a new ElasticsearchRoleLister.
func NewElasticsearchRoleLister(indexer cache.Indexer) ElasticsearchRoleLister {
	return &elasticsearchRoleLister{indexer: indexer}
}

// List lists all ElasticsearchRoles in the indexer.
func (s *elasticsearchRoleLister) List(selector labels.Selector) (ret []*v1alpha1.ElasticsearchRole, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ElasticsearchRole))
	})
	return ret, err
}

// ElasticsearchRoles returns an object that can list and get ElasticsearchRoles.
func (s *elasticsearchRoleLister) ElasticsearchRoles(namespace string) ElasticsearchRoleNamespaceLister {
	return elasticsearchRoleNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ElasticsearchRoleNamespaceLister helps list and get ElasticsearchRoles.
// All objects returned here must be treated as read-only.
type ElasticsearchRoleNamespaceLister interface {
	// List lists all ElasticsearchRoles in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ElasticsearchRole, err error)
	// Get retrieves the ElasticsearchRole from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.ElasticsearchRole, error)
	ElasticsearchRoleNamespaceListerExpansion
}

// elasticsearchRoleNamespaceLister implements the ElasticsearchRoleNamespaceLister
// interface.
type elasticsearchRoleNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ElasticsearchRoles in the indexer for a given namespace.
func (s elasticsearchRoleNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.ElasticsearchRole, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ElasticsearchRole))
	})
	return ret, err
}

// Get retrieves the ElasticsearchRole from the indexer for a given namespace and name.
func (s elasticsearchRoleNamespaceLister) Get(name string) (*v1alpha1.ElasticsearchRole, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("elasticsearchrole"), name)
	}
	return obj.(*v1alpha1.ElasticsearchRole), nil
}
