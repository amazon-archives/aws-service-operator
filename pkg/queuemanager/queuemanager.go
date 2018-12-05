package queuemanager

import (
	"encoding/json"
	"strings"

	"github.com/awslabs/aws-service-operator/pkg/config"
	"github.com/awslabs/aws-service-operator/pkg/helpers"
)

// New will return the QueueManager
func New() *QueueManager {
	return &QueueManager{
		handlers: make(map[string]Handler),
	}
}

// Get will return the handler func
func (q *QueueManager) Get(name string) (handler Handler, ok bool) {
	q.lock.RLock()
	defer q.lock.RUnlock()
	if handler, ok = q.handlers[name]; ok {
		return handler, ok
	}
	return handler, false
}

// Add will add a new handler func
func (q *QueueManager) Add(name string, handlerFunc Handler) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.handlers[name] = handlerFunc
}

// Keys will return the list of topic ARNs
func (q *QueueManager) Keys() []string {
	q.lock.RLock()
	defer q.lock.RUnlock()
	keys := []string{}
	for key, _ := range q.handlers {
		keys = append(keys, key)
	}
	return keys
}

// HandleMessage will stub the handler for processing messages
func (f HandlerFunc) HandleMessage(config config.Config, msg *MessageBody) error {
	return f(config, msg)
}

// ParseMessage will take the message attribute and make it readable
func (m *MessageBody) ParseMessage() error {
	m.Updatable = false
	resp := make(map[string]string)
	items := strings.Split(m.Message, "\n")
	for _, item := range items {
		x := strings.Split(item, "=")
		key := x[0]
		if key != "" {
			s := x[1]
			s = s[1 : len(s)-1]
			resp[key] = s
		}
	}
	m.ParsedMessage = resp

	var resourceProperties ResourceProperties
	if resp["ResourceProperties"] != "null" {
		err := json.Unmarshal([]byte(resp["ResourceProperties"]), &resourceProperties)
		if err != nil {
			return err
		}
		m.ResourceProperties = resourceProperties
		for _, tag := range resourceProperties.Tags {
			switch tag.Key {
			case "Namespace":
				m.Namespace = tag.Value
			case "ResourceName":
				m.ResourceName = tag.Value
			}
		}
		if m.Namespace != "" && m.ResourceName != "" {
			m.Updatable = true
		}
	}
	return nil
}

// IsComplete returns a simple status instead of the raw CFT resp
func (m *MessageBody) IsComplete() bool {
	return helpers.IsStackComplete(m.ParsedMessage["ResourceStatus"], true)
}
