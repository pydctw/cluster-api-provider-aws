/*
Copyright 2021 The Kubernetes Authors.

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

package v1alpha4

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCreateAWSClusterStaticIdentityValidation(t *testing.T) {
	tests := []struct {
		name      string
		selectors map[string]string
		wantError bool
	}{
		{
			name:      "should not return error for valid selector",
			selectors: map[string]string{"foo": "bar"},
			wantError: false,
		},
		{
			name:      "should return error for invalid selector",
			selectors: map[string]string{"-123-foo": "bar"},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identity := &AWSClusterStaticIdentity{
				TypeMeta: metav1.TypeMeta{
					APIVersion: GroupVersion.String(),
					Kind:       string(ClusterStaticIdentityKind),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "static",
				},
				Spec: AWSClusterStaticIdentitySpec{
					AWSClusterIdentitySpec: AWSClusterIdentitySpec{
						AllowedNamespaces: &AllowedNamespaces{
							Selector: metav1.LabelSelector{
								MatchLabels: tt.selectors,
							},
						},
					},
					SecretRef: "test-secret",
				},
			}

			ctx := context.TODO()
			if err := testEnv.Create(ctx, identity); (err != nil) != tt.wantError {
				t.Errorf("ValidateCreate() error = %v, wantErr %v", err, tt.wantError)
			}
			testEnv.Delete(ctx, identity)
		})
	}
}

func TestAWSClusterStaticValidateUpdate(t *testing.T) {
	staticIdentity := &AWSClusterStaticIdentity{
		TypeMeta: metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       string(ClusterStaticIdentityKind),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "static",
		},
		Spec: AWSClusterStaticIdentitySpec{
			SecretRef: "test-secret",
		},
	}

	ctx := context.TODO()
	defer testEnv.Delete(ctx, staticIdentity)

	if err := testEnv.Create(ctx, staticIdentity); err != nil {
		t.Errorf("staticIdentity creation failed %v", err)
	}

	tests := []struct {
		name      string
		identity  *AWSClusterStaticIdentity
		wantError bool
	}{
		{
			name: "do not allow any spec changes",
			identity: &AWSClusterStaticIdentity{
				Spec: AWSClusterStaticIdentitySpec{
					SecretRef:"test",
				},
			},
			wantError: true,
		},
		{
			name:      "no error when updating the same object",
			identity:  staticIdentity,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identity := tt.identity.DeepCopy()
			identity.TypeMeta = metav1.TypeMeta{
				APIVersion: GroupVersion.String(),
				Kind:       string(ClusterStaticIdentityKind),
			}
			ctx := context.TODO()
			if err := testEnv.Update(ctx, identity); (err != nil) != tt.wantError {
				t.Errorf("ValidateUpdate() error = %v, wantErr %v", err, tt.wantError)
			}
		})
	}
}

func TestAWSClusterStaticIdentityUpdateLabelSelectorValidation(t *testing.T) {
	staticIdentity := &AWSClusterStaticIdentity{
		TypeMeta: metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       string(ClusterStaticIdentityKind),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "static",
		},
		Spec: AWSClusterStaticIdentitySpec{
			AWSClusterIdentitySpec: AWSClusterIdentitySpec{
				AllowedNamespaces: &AllowedNamespaces{
					Selector: metav1.LabelSelector{
						MatchLabels: map[string]string{"foo": "bar"},
					},
				},
			},
			SecretRef: "test-secret",
		},
	}

	ctx := context.TODO()
	defer testEnv.Delete(ctx, staticIdentity)

	if err := testEnv.Create(ctx, staticIdentity); err != nil {
		t.Errorf("staticIdentity creation failed %v", err)
	}

	tests := []struct {
		name      string
		identity  *AWSClusterStaticIdentity
		wantError bool
	}{
		{
			name:      "should not return error for valid selector",
			identity:  staticIdentity,
			wantError: false,
		},
		{
			name: "should return error for invalid selector",
			identity: &AWSClusterStaticIdentity{
				ObjectMeta: metav1.ObjectMeta{
					Name: "static",
				},
				Spec: AWSClusterStaticIdentitySpec{
					AWSClusterIdentitySpec: AWSClusterIdentitySpec{
						AllowedNamespaces: &AllowedNamespaces{
							Selector: metav1.LabelSelector{
								MatchLabels: map[string]string{"-foo-123": "bar"},
							},
						},
					},
				},
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identity := tt.identity.DeepCopy()
			ctx := context.TODO()
			if err := testEnv.Update(ctx, identity); (err != nil) != tt.wantError {
				t.Errorf("ValidateUpdate() error = %v, wantErr %v", err, tt.wantError)
			}
		})
	}
}
