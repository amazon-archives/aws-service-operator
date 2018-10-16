package operator

import (
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type operator struct {
	plural    string
	namespace string
	handlers  cache.ResourceEventHandlerFuncs
	client    rest.Interface
}

func New(plural, namespace string, handlers cache.ResourceEventHandlerFuncs, client rest.Interface) *operator {
	return &operator{
		plural:    plural,
		namespace: namespace,
		handlers:  handlers,
		client:    client,
	}
}

func (o *operator) Watch(obj runtime.Object, done <-chan struct{}) error {
	source := cache.NewListWatchFromClient(o.client, o.plural, o.namespace, fields.Everything())
	_, controller := cache.NewInformer(source, obj, 5, o.handlers)

	go controller.Run(done)
	<-done
	return nil
}
