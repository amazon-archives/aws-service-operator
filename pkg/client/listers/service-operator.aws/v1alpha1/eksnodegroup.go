/*
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
	v1alpha1 "github.com/awslabs/aws-service-operator/pkg/apis/service-operator.aws/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// EKSNodeGroupLister helps list EKSNodeGroups.
type EKSNodeGroupLister interface {
	// List lists all EKSNodeGroups in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.EKSNodeGroup, err error)
	// EKSNodeGroups returns an object that can list and get EKSNodeGroups.
	EKSNodeGroups(namespace string) EKSNodeGroupNamespaceLister
	EKSNodeGroupListerExpansion
}

// eKSNodeGroupLister implements the EKSNodeGroupLister interface.
type eKSNodeGroupLister struct {
	indexer cache.Indexer
}

// NewEKSNodeGroupLister returns a new EKSNodeGroupLister.
func NewEKSNodeGroupLister(indexer cache.Indexer) EKSNodeGroupLister {
	return &eKSNodeGroupLister{indexer: indexer}
}

// List lists all EKSNodeGroups in the indexer.
func (s *eKSNodeGroupLister) List(selector labels.Selector) (ret []*v1alpha1.EKSNodeGroup, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.EKSNodeGroup))
	})
	return ret, err
}

// EKSNodeGroups returns an object that can list and get EKSNodeGroups.
func (s *eKSNodeGroupLister) EKSNodeGroups(namespace string) EKSNodeGroupNamespaceLister {
	return eKSNodeGroupNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// EKSNodeGroupNamespaceLister helps list and get EKSNodeGroups.
type EKSNodeGroupNamespaceLister interface {
	// List lists all EKSNodeGroups in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.EKSNodeGroup, err error)
	// Get retrieves the EKSNodeGroup from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.EKSNodeGroup, error)
	EKSNodeGroupNamespaceListerExpansion
}

// eKSNodeGroupNamespaceLister implements the EKSNodeGroupNamespaceLister
// interface.
type eKSNodeGroupNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all EKSNodeGroups in the indexer for a given namespace.
func (s eKSNodeGroupNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.EKSNodeGroup, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.EKSNodeGroup))
	})
	return ret, err
}

// Get retrieves the EKSNodeGroup from the indexer for a given namespace and name.
func (s eKSNodeGroupNamespaceLister) Get(name string) (*v1alpha1.EKSNodeGroup, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("eksnodegroup"), name)
	}
	return obj.(*v1alpha1.EKSNodeGroup), nil
}
