package helpers

import (
	"github.com/awslabs/aws-service-operator/pkg/config"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateConfigMap will create a Kubernetes Servic Using ExternalName types
func CreateConfigMap(config config.Config, resource interface{}, svcName string, svcNamespace string, configMapTemplate map[string]string) string {
	logger := config.Logger
	cmData := map[string]string{}
	for key, value := range configMapTemplate {
		tempValue, err := Templatize(value, Data{Obj: resource, Config: config, Helpers: New()})
		if err != nil {
			logger.WithError(err).Error("error parsing config map template")
			return ""
		}
		cmData[key] = tempValue
	}

	configMap := &apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: KubernetesResourceName(svcName),
		},
		Data: cmData,
	}

	newConfigMap, err := config.KubeClientset.CoreV1().ConfigMaps(svcNamespace).Create(configMap)
	if err != nil {
		logger.WithError(err).Error("error creating configmap")
	}
	return newConfigMap.Name
}
