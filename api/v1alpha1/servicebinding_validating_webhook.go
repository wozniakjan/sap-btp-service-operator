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

package v1alpha1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var servicebindinglog = logf.Log.WithName("servicebinding-resource")

func (sb *ServiceBinding) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(sb).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-services-cloud-sap-com-v1alpha1-servicebinding,mutating=false,failurePolicy=fail,groups=services.cloud.sap.com,resources=servicebindings,versions=v1alpha1,name=vservicebinding.kb.io,sideEffects=None,admissionReviewVersions=v1beta1;v1

var _ webhook.Validator = &ServiceBinding{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (sb *ServiceBinding) ValidateCreate() error {
	servicebindinglog.Info("validate create", "name", sb.Name)
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (sb *ServiceBinding) ValidateUpdate(old runtime.Object) error {
	servicebindinglog.Info("validate update", "name", sb.Name)

	if sb.specChanged(old) && sb.Status.BindingID != "" {
		return fmt.Errorf("updating service bindings is not supported")
	}

	return nil
}

func (sb *ServiceBinding) specChanged(old runtime.Object) bool {
	oldBinding := old.(*ServiceBinding)

	if changed := sb.paramsFromChanged(oldBinding); changed {
		return true
	}

	return sb.Spec.ExternalName != oldBinding.Spec.ExternalName ||
		sb.Spec.ServiceInstanceName != oldBinding.Spec.ServiceInstanceName ||
		// TODO + labels
		//r.Spec.Labels != oldBinding.Spec.Labels ||
		sb.Spec.Parameters.String() != oldBinding.Spec.Parameters.String() ||
		sb.Spec.SecretName != oldBinding.Spec.SecretName
}

func (sb *ServiceBinding) paramsFromChanged(oldBinding *ServiceBinding) bool {
	if len(sb.Spec.ParametersFrom) != len(oldBinding.Spec.ParametersFrom) {
		return true
	}
	for i, paramFrom := range sb.Spec.ParametersFrom {
		if paramFrom.SecretKeyRef != nil && oldBinding.Spec.ParametersFrom[i].SecretKeyRef != nil {
			if *paramFrom.SecretKeyRef != *oldBinding.Spec.ParametersFrom[i].SecretKeyRef {
				return true
			}
		} else if paramFrom.SecretKeyRef != oldBinding.Spec.ParametersFrom[i].SecretKeyRef {
			return true
		}
	}
	return false
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (sb *ServiceBinding) ValidateDelete() error {
	servicebindinglog.Info("validate delete", "name", sb.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
