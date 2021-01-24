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

package database

import (
	"kubevault.dev/operator/pkg/vault/role"

	vaultapi "github.com/hashicorp/vault/api"
	rbacv1 "k8s.io/api/rbac/v1"
)

type DatabaseRoleInterface interface {
	role.RoleInterface

	// EnableDatabase enables database secret engine
	EnableDatabase() error

	// IsDatabaseEnabled checks whether database is enabled or not
	IsDatabaseEnabled() (bool, error)

	// DeleteRole deletes role
	DeleteRole(name string) (int, error)
}

type DatabaseCredentialManager interface {
	// Gets credential from vault
	GetCredential() (*vaultapi.Secret, error)

	// Creates a kubernetes secret containing postgres credential
	CreateSecret(name string, namespace string, credential *vaultapi.Secret) error

	// Creates kubernetes role
	CreateRole(name string, namespace string, secretName string) error

	// Creates kubernetes role binding
	CreateRoleBinding(name string, namespace string, roleName string, subjects []rbacv1.Subject) error

	IsLeaseExpired(leaseID string) (bool, error)

	RevokeLease(leaseID string) error
}
