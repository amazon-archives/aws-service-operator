package helpers

import (
	"strconv"
	"strings"

	"github.com/awslabs/aws-service-operator/pkg/config"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateExternalNameService will create a Kubernetes Servic Using ExternalName types
func CreateExternalNameService(config config.Config, resource interface{}, svcName string, svcNamespace string, externalNameTemplate string, svcPortTemplate string) string {
	logger := config.Logger

	externalName, err := Templatize(externalNameTemplate, Data{Obj: resource, Config: config})
	if err != nil {
		logger.WithError(err).Error("error parsing external name template")
		return ""
	}

	svcPortString, err := Templatize(svcPortTemplate, Data{Obj: resource, Config: config})
	if err != nil {
		logger.WithError(err).Error("error parsing service port template")
		return ""
	}

	svcPortStripSlash := strings.Replace(svcPortString, "\\", "", -1)

	svcPortInt64, err := strconv.ParseInt(svcPortStripSlash, 0, 16)
	if err != nil {
		logger.WithError(err).Error("error converting service port template string to int")
		return ""
	}

	// ParseInt only returns an int64, must convert to int32 for apiv1.ServicePort field
	svcPort := int32(svcPortInt64)

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: KubernetesResourceName(svcName),
		},
		Spec: apiv1.ServiceSpec{
			Type:         apiv1.ServiceTypeExternalName,
			ExternalName: externalName,
			Ports: []apiv1.ServicePort{
				apiv1.ServicePort{
					Port: svcPort,
				},
			},
		},
	}

	newService, err := config.KubeClientset.CoreV1().Services(svcNamespace).Create(service)
	if err != nil {
		logger.WithError(err).Error("error creating service")
	}
	return newService.Name
}

// func serviceType(str string) apiv1.ServiceType {
// 	var sType apiv1.ServiceType
// 	switch str {
// 	case "ClusterIP":
// 		sType = apiv1.ServiceTypeClusterIP
// 	case "NodePort":
// 		sType = apiv1.ServiceTypeNodePort
// 	case "LoadBalancer":
// 		sType = apiv1.ServiceTypeLoadBalancer
// 	case "ExternalName":
// 		sType = apiv1.ServiceTypeExternalName
// 	default:
// 		sType = apiv1.ServiceTypeClusterIP
// 	}
// 	return sType
// }
