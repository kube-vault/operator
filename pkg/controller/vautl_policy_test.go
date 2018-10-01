package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/appscode/pat"
	policyapi "github.com/kubevault/operator/apis/policy/v1alpha1"
	csfake "github.com/kubevault/operator/client/clientset/versioned/fake"
	"github.com/kubevault/operator/pkg/vault/policy"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kfake "k8s.io/client-go/kubernetes/fake"
)

type fakePolicy struct {
	errInPutPolicy bool
}

func (f *fakePolicy) EnsurePolicy(n, p string) error {
	if f.errInPutPolicy {
		return errors.New("error")
	}
	return nil
}

func (f *fakePolicy) DeletePolicy(n string) error {
	return nil
}

func simpleVaultPolicy() *policyapi.VaultPolicy {
	return &policyapi.VaultPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "simple",
			Namespace:  "test",
			Finalizers: []string{VaultPolicyFinalizer},
		},
		Spec: policyapi.VaultPolicySpec{
			Policy: "simple {}",
		},
	}
}

func validVaultPolicy(vAddr, tokenSecret string) *policyapi.VaultPolicy {
	return &policyapi.VaultPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "ok",
			Namespace:  "test",
			Finalizers: []string{VaultPolicyFinalizer},
		},
		Spec: policyapi.VaultPolicySpec{
			Policy: "simple {}",
			Vault: &policyapi.Vault{
				Address:             vAddr,
				TokenSecret:         tokenSecret,
				SkipTLSVerification: true,
			},
		},
	}
}

func vaultTokenSecret() *core.Secret {
	return &core.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "vault",
			Namespace: "test",
		},
		Data: map[string][]byte{
			"token": []byte("root"),
		},
	}
}

func NewFakeVaultServer() *httptest.Server {
	m := pat.New()
	m.Del("/v1/sys/policy/ok", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	m.Del("/v1/sys/policy/simple", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	return httptest.NewServer(m)
}

func TestReconcilePolicy(t *testing.T) {
	cases := []struct {
		testName     string
		vPolicy      *policyapi.VaultPolicy
		pClient      policy.Policy
		expectStatus string
		expectErr    bool
	}{
		{
			testName:     "reconcile successful",
			vPolicy:      simpleVaultPolicy(),
			pClient:      &fakePolicy{},
			expectStatus: string(policyapi.PolicySuccess),
			expectErr:    false,
		},
		{
			testName:     "reconcile unsuccessful, error occure in EnsurePolicy",
			vPolicy:      simpleVaultPolicy(),
			pClient:      &fakePolicy{errInPutPolicy: true},
			expectStatus: string(policyapi.PolicyFailed),
			expectErr:    true,
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			ctrl := &VaultController{
				extClient: csfake.NewSimpleClientset(simpleVaultPolicy()),
			}

			err := ctrl.reconcilePolicy(c.vPolicy, c.pClient)
			if c.expectErr {
				assert.NotNil(t, err, "expected error")
			} else {
				assert.Nil(t, err)
			}
			if c.expectStatus != "" {
				p, err := ctrl.extClient.PolicyV1alpha1().VaultPolicies(c.vPolicy.Namespace).Get(c.vPolicy.Name, metav1.GetOptions{})
				if assert.Nil(t, err) {
					assert.Condition(t, func() (success bool) {
						return c.expectStatus == string(p.Status.Status)
					}, ".status.status should match")
				}
			}
		})
	}
}

func TestFinalizePolicy(t *testing.T) {
	srv := NewFakeVaultServer()
	defer srv.Close()

	kc := kfake.NewSimpleClientset(vaultTokenSecret())

	cases := []struct {
		testName  string
		vPolicy   *policyapi.VaultPolicy
		expectErr bool
	}{
		{
			testName:  "no error, valid VaultPolicy",
			vPolicy:   validVaultPolicy(srv.URL, vaultTokenSecret().Name),
			expectErr: false,
		},
		{
			testName:  "no error, VaultPolicy doesn't exist",
			vPolicy:   nil,
			expectErr: false,
		},
		{
			testName:  "error, invalid VaultPolicy",
			vPolicy:   validVaultPolicy("invalid", vaultTokenSecret().Name),
			expectErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			pc := csfake.NewSimpleClientset().PolicyV1alpha1()
			if c.vPolicy != nil {
				_, err := pc.VaultPolicies(c.vPolicy.Namespace).Create(c.vPolicy)
				assert.Nil(t, err)
			} else {
				c.vPolicy = simpleVaultPolicy()
			}

			err := finalizePolicy(pc, kc, c.vPolicy)
			if c.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestRunPolicyFinalizer(t *testing.T) {
	srv := NewFakeVaultServer()
	defer srv.Close()
	ctrl := &VaultController{
		extClient:     csfake.NewSimpleClientset(simpleVaultPolicy(), validVaultPolicy(srv.URL, vaultTokenSecret().Name)),
		kubeClient:    kfake.NewSimpleClientset(vaultTokenSecret()),
		finalizerInfo: NewMapFinalizer(),
	}
	ctrl.finalizerInfo.Add(simpleVaultPolicy().GetKey())

	cases := []struct {
		testName  string
		vPolicy   *policyapi.VaultPolicy
		completed bool
	}{
		{
			testName:  "remove finalizer successfully, valid VaultPolicy",
			vPolicy:   validVaultPolicy(srv.URL, vaultTokenSecret().Name),
			completed: true,
		},
		{
			testName:  "remove finalizer successfully, invalid VaultPolicy",
			vPolicy:   validVaultPolicy("invalid", vaultTokenSecret().Name),
			completed: true,
		},
		{
			testName:  "already processing finalizer",
			vPolicy:   simpleVaultPolicy(),
			completed: false,
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			ctrl.runPolicyFinalizer(c.vPolicy, 3*time.Second, 1*time.Second)
			if c.completed {
				assert.Condition(t, func() (success bool) {
					return !ctrl.finalizerInfo.IsAlreadyProcessing(c.vPolicy.GetKey())
				}, "IsAlreadyProcessing(key) should be false")

			} else {
				assert.Condition(t, func() (success bool) {
					return ctrl.finalizerInfo.IsAlreadyProcessing(c.vPolicy.GetKey())
				}, "IsAlreadyProcessing(key) should be true")
			}
		})
	}
}
