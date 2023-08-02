/*
Copyright 2023.

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

package v1

import (
	"fmt"
	"regexp"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var oslog = logf.Log.WithName("os-resource")

func (r *OS) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-apps-my-domain-v1-os,mutating=true,failurePolicy=fail,sideEffects=None,groups=apps.my.domain,resources=os,verbs=create;update,versions=v1,name=mos.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &OS{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *OS) Default() {
	oslog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-apps-my-domain-v1-os,mutating=false,failurePolicy=fail,sideEffects=None,groups=apps.my.domain,resources=os,verbs=create;update,versions=v1,name=vos.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &OS{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *OS) ValidateCreate() error {
	oslog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return r.validateOS()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *OS) ValidateUpdate(old runtime.Object) error {
	oslog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return r.validateOS()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *OS) ValidateDelete() error {
	oslog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (r *OS) validateOS() error {
	var validImagetypes = []string{"docker", "containerd", "disk"}
	var validOpstypes = []string{"upgrade", "config", "rollback"}

	if !contains(validImagetypes, r.Spec.ImageType) {
		return fmt.Errorf("invalid imagetype: %s", r.Spec.ImageType)
	}
	if !contains(validOpstypes, r.Spec.OpsType) {
		return fmt.Errorf("invalid opstype: %s", r.Spec.OpsType)
	}
	if r.Spec.OSVersion == "" {
		return fmt.Errorf("osversion is required")
	}
	if r.Spec.MaxUnavailable < 0 {
		return fmt.Errorf("maxunavailable should be greater than or equal to 0")
	}
	if !isContainerImage(r.Spec.ContainerImage) {
		return fmt.Errorf("invalid containerimage: %s", r.Spec.ContainerImage)
	}
	if r.Spec.ImageType == "disk" {
		if !isHTTPS(r.Spec.ImageURL) && !isHTTP(r.Spec.ImageURL) {
			return fmt.Errorf("only http or https protocols are supported : %s", r.Spec.ImageURL)
		}
		if r.Spec.CheckSum == "" {
			return fmt.Errorf("checksum is required")
		}
		if isHTTP(r.Spec.ImageURL) {
			if !r.Spec.FlagSafe {
				return fmt.Errorf("flagsafe should be true when using http protocol")
			}
		}
		if isHTTPS(r.Spec.ImageURL) {
			if r.Spec.CaCert == "" {
				return fmt.Errorf("cacert is required")
			}
			if r.Spec.MTLS {
				if r.Spec.ClientCert == "" {
					return fmt.Errorf("clientcert is required")
				}
				if r.Spec.ClientKey == "" {
					return fmt.Errorf("clientkey is required")
				}
			}
		}
	}
	return nil
}

// 定义 contains 函数，用于判断一个字符串是否包含在一个字符串数组中
func contains(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

// 定义 isHTTP 函数，用于判断一个 URL 是否使用了 HTTP 协议
func isHTTP(url string) bool {
	return len(url) >= 7 && url[:8] == "http://"
}

// 定义 isHTTPS 函数，用于判断一个 URL 是否使用了 HTTPS 协议
func isHTTPS(url string) bool {
	return len(url) >= 8 && url[:8] == "https://"
}

// 定义 isMap 函数，用于判断一个 interface{} 是否为 map 类型
func isMap(v interface{}) bool {
	return v != nil && fmt.Sprintf("%T", v) == "map[string]interface {}"
}

func isContainerImage(str string) bool {
	// 定义容器镜像正则表达式
	pattern := `^([\w.-]+\/)?[\w.-]+(:[\w.-]+)?$`

	// 编译正则表达式
	regex, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Invalid regex:", err)
		return false
	}

	// 匹配字符串
	return regex.MatchString(str)
}
