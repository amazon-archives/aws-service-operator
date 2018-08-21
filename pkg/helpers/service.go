package helpers

import (
	"github.com/christopherhein/aws-operator/pkg/config"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Data wrapps the object that is needed for the services
type Data struct {
	Obj    interface{}
	Config *config.Config
}

// CreateExternalNameService will create a Kubernetes Servic Using ExternalName types
func CreateExternalNameService(config *config.Config, resource interface{}, svcName string, svcNamespace string, externalNameTemplate string, svcPort int32) string {
	logger := config.Logger

	externalName, err := Templatize(externalNameTemplate, Data{Obj: resource, Config: config})
	if err != nil {
		logger.WithError(err).Error("error parsing external name template")
		return ""
	}

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: svcName,
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

	newService, err := config.Context.Clientset.CoreV1().Services(svcNamespace).Create(service)
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
