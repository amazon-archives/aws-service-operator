/*
Copyright 2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License"). You may
not use this file except in compliance with the License. A copy of the
License is located at

     http://aws.amazon.com/apache2.0/

or in the "license" file accompanying this file. This file is distributed
on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
express or implied. See the License for the specific language governing
permissions and limitations under the License.
*/

package cloudformation

import (
	"reflect"
	"strconv"
	"strings"

	"k8s.io/klog"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
)

const (
	tagName          = "cloudformation"
	parameterKeyName = "Parameter"
	outputKeyName    = "Output"
)

type encoderFunc func(types map[string]string, field reflect.Type, index int, v reflect.Value, typeKey string)

// MarshalParameters will create Cloudformation parameters for each of the
// outputs in a go struct
func MarshalParameters(v interface{}) []*cfn.Parameter {
	out := []*cfn.Parameter{}
	types := map[string]string{}

	MarshalTypes(types, v, parameterKeyName)
	for key, val := range types {
		param := &cfn.Parameter{}
		param.SetParameterKey(key)
		param.SetParameterValue(val)
		out = append(out, param)
	}

	return out
}

// MarshalOutputs will create Cloudformation outputs for each of the outputs in
// a go struct
func MarshalOutputs(v interface{}) []*cfn.Output {
	out := []*cfn.Output{}
	types := map[string]string{}

	MarshalTypes(types, v, outputKeyName)
	for key, val := range types {
		output := &cfn.Output{}
		output.SetOutputKey(key)
		output.SetOutputValue(val)
		out = append(out, output)
	}

	return out
}

// MarshalTypes will generate a map of strings by the specific key type
func MarshalTypes(types map[string]string, v interface{}, typeKey string) {
	value := reflect.ValueOf(v)
	translateRecursive(types, value, typeKey, "", "")
	return
}

func translateRecursive(types map[string]string, value reflect.Value, typeKey, tagKey, tagType string) {
	switch value.Kind() {
	case reflect.Ptr:
		vValue := value.Elem()
		if !vValue.IsValid() {
			return
		}
		translateRecursive(types, vValue, typeKey, tagKey, tagType)
	case reflect.Interface:
		vValue := value.Elem()
		translateRecursive(types, vValue, typeKey, tagKey, tagType)
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			tagVal := value.Type().Field(i).Tag.Get(tagName)

			if tagVal != "" {
				tagParser := strings.Split(tagVal, ",")
				if tagParser[1] == typeKey {
					tagType = tagParser[1]
					tagKey = tagParser[0]
				} else {
					tagType = ""
					tagKey = ""
				}
			}
			translateRecursive(types, value.Field(i), typeKey, tagKey, tagType)
		}
	case reflect.Slice:
		if value.Len() == 0 && tagKey != "" {
			types[tagKey] = ""
		}
		for i := 0; i < value.Len(); i++ {
			translateRecursive(types, value.Index(i), typeKey, tagKey, tagType)
		}
	case reflect.Map:
		for _, key := range value.MapKeys() {
			vValue := value.MapIndex(key)
			translateRecursive(types, vValue, typeKey, tagKey, tagType)
		}
	case reflect.Bool:
		if tagKey != "" {
			types[tagKey] = strconv.FormatBool(value.Bool())
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if tagKey != "" {
			types[tagKey] = strconv.FormatInt(value.Int(), 10)
		}
	case reflect.String:
		if tagKey != "" {
			tagName, ok := types[tagKey]
			strValue := value.String()
			if !ok {
				types[tagKey] = strValue
			} else {
				if strValue != "" {
					stringSlice := []string{tagName, strValue}
					types[tagKey] = strings.Join(stringSlice, ", ")
				}
			}
		}
	default:
		klog.V(3).Info("parsed unsupported struct type")
	}
}
