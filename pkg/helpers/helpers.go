package helpers

import (
	"bytes"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	awsclient "github.com/awslabs/aws-service-operator/pkg/client/clientset/versioned/typed/service-operator.aws/v1alpha1"
	"github.com/awslabs/aws-service-operator/pkg/config"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KubernetesResourceName returns the resource name for other components
func KubernetesResourceName(name string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9_-]+")
	return reg.ReplaceAllString(name, "-")
}

// StackName will return the proper stack name for each component
func StackName(clusterName string, resourceType string, name string, namespace string) string {
	nameParts := []string{clusterName, resourceType, name, namespace}
	return KubernetesResourceName(strings.Join(nameParts, "-"))
}

// Stringify will create a string based on the params
func Stringify(attr interface{}) string {
	switch reflect.TypeOf(attr).Name() {
	case "bool":
		return strconv.FormatBool(attr.(bool))
	case "string":
		return attr.(string)
	case "int":
		return strconv.Itoa(attr.(int))
	}
	return ""
}

// CreateParam returns a new prefilled cloudformation param
func CreateParam(key string, value string) *cloudformation.Parameter {
	param := &cloudformation.Parameter{}
	param.SetParameterKey(key)
	param.SetParameterValue(value)
	return param
}

// CreateTag returns a new prefilled cloudformation tag
func CreateTag(key string, value string) *cloudformation.Tag {
	tag := &cloudformation.Tag{}
	tag.SetKey(key)
	tag.SetValue(value)
	return tag
}

// IsStackComplete will determine if it's in a state to process
func IsStackComplete(status string, defaultRet bool) bool {
	switch status {
	case "CREATE_COMPLETE":
		return true
	case "UPDATE_COMPLETE":
		return true
	case "DELETE_COMPLETE":
		return false
	case "ROLLBACK_COMPLETE":
		return false
	}
	return defaultRet
}

// Templatize returns the proper values based on the templating
func Templatize(tempStr string, data interface{}) (resp string, err error) {
	t := template.New("templating")
	t, err = t.Parse(string(tempStr))
	if err != nil {
		return
	}

	var tpl bytes.Buffer
	err = t.Execute(&tpl, data)
	return tpl.String(), err
}

// GetCloudFormationTemplate will return the url to the CFT from the CFT resource
func GetCloudFormationTemplate(config config.Config, rType string, name string, namespace string) string {
	logger := config.Logger
	clientSet, _ := awsclient.NewForConfig(config.RESTConfig)

	var cName string
	var cNamespace string
	if name == "" {
		cName = rType
	}
	if namespace == "" {
		cNamespace = config.DefaultNamespace
	}

	resource, err := clientSet.CloudFormationTemplates(cNamespace).Get(cName, metav1.GetOptions{})
	if err != nil {
		logger.WithError(err).Error("error getting cloudformation template returning fallback template")
		return "https://s3-us-west-2.amazonaws.com/cloudkit-templates/" + rType + ".yaml"
	}

	logger.WithFields(logrus.Fields{
		"namespace": cNamespace,
		"name":      cName,
		"url":       resource.Output.URL,
	}).Info("found cloudformation template")
	return resource.Output.URL
}
