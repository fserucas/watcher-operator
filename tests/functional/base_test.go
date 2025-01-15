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

package functional

import (
	"fmt"

	. "github.com/onsi/gomega" //revive:disable:dot-imports

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	condition "github.com/openstack-k8s-operators/lib-common/modules/common/condition"
	watcherv1 "github.com/openstack-k8s-operators/watcher-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

func GetDefaultWatcherSpec() map[string]interface{} {
	return map[string]interface{}{
		"databaseInstance": "openstack",
		"secret":           SecretName,
	}
}

// Second Watcher Spec to test proper parameters substitution
func GetNonDefaultWatcherSpec() map[string]interface{} {
	return map[string]interface{}{
		"apiContainerImageURL": "fake-API-Container-URL",
		"secret":               SecretName,
		"preserveJobs":         true,
		"databaseInstance":     "fakeopenstack",
		"serviceUser":          "fakeuser",
		"customServiceConfig":  "# Global config",
		"apiServiceTemplate": map[string]interface{}{
			"replicas":            2,
			"nodeSelector":        map[string]string{"foo": "bar"},
			"customServiceConfig": "# Service config",
		},
		"tls": map[string]interface{}{
			"caBundleSecretName": "combined-ca-bundle",
		},
	}
}

func GetDefaultWatcherAPISpec() map[string]interface{} {
	return map[string]interface{}{
		"databaseInstance":  "openstack",
		"secret":            SecretName,
		"memcachedInstance": "memcached",
		"serviceAccount":    "watcher-sa",
		"containerImage":    "test://watcher",
	}
}

func GetDefaultWatcherApplierSpec() map[string]interface{} {
	return map[string]interface{}{
		"databaseInstance":  "openstack",
		"secret":            SecretName,
		"memcachedInstance": "memcached",
		"serviceAccount":    "watcher-sa",
		"containerImage":    "test://watcher",
	}
}

func GetDefaultWatcherDecisionEngineSpec() map[string]interface{} {
	return map[string]interface{}{
		"databaseInstance":  "openstack",
		"secret":            SecretName,
		"memcachedInstance": "memcached",
		"serviceAccount":    "watcher-sa",
		"containerImage":    "test://watcher",
	}
}

func CreateWatcher(name types.NamespacedName, spec map[string]interface{}) client.Object {
	raw := map[string]interface{}{
		"apiVersion": "watcher.openstack.org/v1beta1",
		"kind":       "Watcher",
		"metadata": map[string]interface{}{
			"name":      name.Name,
			"namespace": name.Namespace,
		},
		"spec": spec,
	}
	return th.CreateUnstructured(raw)
}

func GetWatcher(name types.NamespacedName) *watcherv1.Watcher {
	instance := &watcherv1.Watcher{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(ctx, name, instance)).Should(Succeed())
	}, timeout, interval).Should(Succeed())
	return instance
}

func WatcherConditionGetter(name types.NamespacedName) condition.Conditions {
	instance := GetWatcher(name)
	return instance.Status.Conditions
}

func CreateWatcherAPI(name types.NamespacedName, spec map[string]interface{}) client.Object {
	raw := map[string]interface{}{
		"apiVersion": "watcher.openstack.org/v1beta1",
		"kind":       "WatcherAPI",
		"metadata": map[string]interface{}{
			"name":      name.Name,
			"namespace": name.Namespace,
		},
		"spec": spec,
	}
	return th.CreateUnstructured(raw)
}

func GetWatcherAPI(name types.NamespacedName) *watcherv1.WatcherAPI {
	instance := &watcherv1.WatcherAPI{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(ctx, name, instance)).Should(Succeed())
	}, timeout, interval).Should(Succeed())
	return instance
}

func WatcherAPIConditionGetter(name types.NamespacedName) condition.Conditions {
	instance := GetWatcherAPI(name)
	return instance.Status.Conditions
}

func CreateWatcherApplier(name types.NamespacedName, spec map[string]interface{}) client.Object {
	raw := map[string]interface{}{
		"apiVersion": "watcher.openstack.org/v1beta1",
		"kind":       "WatcherApplier",
		"metadata": map[string]interface{}{
			"name":      name.Name,
			"namespace": name.Namespace,
		},
		"spec": spec,
	}
	return th.CreateUnstructured(raw)
}

func GetWatcherApplier(name types.NamespacedName) *watcherv1.WatcherApplier {
	instance := &watcherv1.WatcherApplier{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(ctx, name, instance)).Should(Succeed())
	}, timeout, interval).Should(Succeed())
	return instance
}

func WatcherApplierConditionGetter(name types.NamespacedName) condition.Conditions {
	instance := GetWatcherApplier(name)
	return instance.Status.Conditions
}

func WatcherDecisionEngineConditionGetter(name types.NamespacedName) condition.Conditions {
	instance := GetWatcherDecisionEngine(name)
	return instance.Status.Conditions
}

func CreateWatcherMessageBusSecret(namespace string, name string) *corev1.Secret {
	s := th.CreateSecret(
		types.NamespacedName{Namespace: namespace, Name: name},
		map[string][]byte{
			"transport_url": []byte(fmt.Sprintf("rabbit://%s/fake", name)),
		},
	)
	logger.Info("Secret created", "name", name)
	return s
}

func CreateWatcherDecisionEngine(name types.NamespacedName, spec map[string]interface{}) client.Object {
	raw := map[string]interface{}{
		"apiVersion": "watcher.openstack.org/v1beta1",
		"kind":       "WatcherDecisionEngine",
		"metadata": map[string]interface{}{
			"name":      name.Name,
			"namespace": name.Namespace,
		},
		"spec": spec,
	}
	return th.CreateUnstructured(raw)
}

func GetWatcherDecisionEngine(name types.NamespacedName) *watcherv1.WatcherDecisionEngine {
	instance := &watcherv1.WatcherDecisionEngine{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(ctx, name, instance)).Should(Succeed())
	}, timeout, interval).Should(Succeed())
	return instance
}
