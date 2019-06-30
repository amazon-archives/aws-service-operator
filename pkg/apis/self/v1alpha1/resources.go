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

package v1alpha1

import "strings"

// AvailableControllers lists all the resources this controller manager handles
var AvailableControllers = map[string][]string{
	"mq":                     []string{"broker", "configuration", "configurationassociation"},
	"apigateway":             []string{"account", "apikey", "authorizer", "basepathmapping", "clientcertificate", "deployment", "documentationpart", "documentationversion", "domainname", "gatewayresponse", "method", "model", "requestvalidator", "apigatewayresource", "restapi", "stage", "usageplan", "usageplankey", "vpclink"},
	"apigatewayv2":           []string{"api", "apimapping", "authorizer", "deployment", "domainname", "integration", "integrationresponse", "model", "route", "routeresponse", "stage"},
	"appmesh":                []string{"mesh", "route", "virtualnode", "virtualrouter", "virtualservice"},
	"appstream":              []string{"directoryconfig", "fleet", "imagebuilder", "stack", "stackfleetassociation", "stackuserassociation", "user"},
	"appsync":                []string{"apikey", "datasource", "functionconfiguration", "graphqlapi", "graphqlschema", "resolver"},
	"applicationautoscaling": []string{"scalabletarget", "scalingpolicy"},
	"athena":                 []string{"namedquery"},
	"autoscaling":            []string{"autoscalinggroup", "launchconfiguration", "lifecyclehook", "scalingpolicy", "scheduledaction"},
	"autoscalingplans":       []string{"scalingplan"},
	"batch":                  []string{"computeenvironment", "jobdefinition", "jobqueue"},
	"budgets":                []string{"budget"},
	"cdk":                    []string{"metadata"},
	"certificatemanager":     []string{"certificate"},
	"cloud9":                 []string{"environmentec2"},
	"cloudfront":             []string{"cloudfrontoriginaccessidentity", "distribution", "streamingdistribution"},
	"cloudtrail":             []string{"trail"},
	"cloudwatch":             []string{"alarm", "dashboard"},
	"codebuild":              []string{"project"},
	"codecommit":             []string{"repository"},
	"codedeploy":             []string{"application", "deploymentconfig", "deploymentgroup"},
	"codepipeline":           []string{"customactiontype", "pipeline", "webhook"},
	"cognito":                []string{"identitypool", "identitypoolroleattachment", "userpool", "userpoolclient", "userpoolgroup", "userpooluser", "userpoolusertogroupattachment"},
	"config":                 []string{"aggregationauthorization", "configrule", "configurationaggregator", "configurationrecorder", "deliverychannel"},
	"dax":                    []string{"cluster", "parametergroup", "subnetgroup"},
	"dlm":                    []string{"lifecyclepolicy"},
	"dms":                    []string{"certificate", "endpoint", "eventsubscription", "replicationinstance", "replicationsubnetgroup", "replicationtask"},
	"datapipeline":           []string{"pipeline"},
	"directoryservice":       []string{"microsoftad", "simplead"},
	"docdb":                  []string{"dbcluster", "dbclusterparametergroup", "dbinstance", "dbsubnetgroup"},
	"dynamodb":               []string{"table"},
	"ec2":                    []string{"capacityreservation", "customergateway", "dhcpoptions", "ec2fleet", "eip", "eipassociation", "egressonlyinternetgateway", "flowlog", "host", "instance", "internetgateway", "launchtemplate", "natgateway", "networkacl", "networkaclentry", "networkinterface", "networkinterfaceattachment", "networkinterfacepermission", "placementgroup", "route", "routetable", "securitygroup", "securitygroupegress", "securitygroupingress", "spotfleet", "subnet", "subnetcidrblock", "subnetnetworkaclassociation", "subnetroutetableassociation", "transitgateway", "transitgatewayattachment", "transitgatewayroute", "transitgatewayroutetable", "transitgatewayroutetableassociation", "transitgatewayroutetablepropagation", "vpc", "vpccidrblock", "vpcdhcpoptionsassociation", "vpcendpoint", "vpcendpointconnectionnotification", "vpcendpointservice", "vpcendpointservicepermissions", "vpcgatewayattachment", "vpcpeeringconnection", "vpnconnection", "vpnconnectionroute", "vpngateway", "vpngatewayroutepropagation", "volume", "volumeattachment"},
	"ecr":                    []string{"repository"},
	"ecs":                    []string{"cluster", "service", "taskdefinition"},
	"efs":                    []string{"filesystem", "mounttarget"},
	"eks":                    []string{"cluster"},
	"emr":                    []string{"cluster", "instancefleetconfig", "instancegroupconfig", "securityconfiguration", "step"},
	"elasticache":            []string{"cachecluster", "parametergroup", "replicationgroup", "securitygroup", "securitygroupingress", "subnetgroup"},
	"elasticbeanstalk":       []string{"application", "applicationversion", "configurationtemplate", "environment"},
	"elasticloadbalancing":   []string{"loadbalancer"},
	"elasticloadbalancingv2": []string{"listener", "listenercertificate", "listenerrule", "loadbalancer", "targetgroup"},
	"elasticsearch":          []string{"domain"},
	"events":                 []string{"eventbuspolicy", "rule"},
	"fsx":                    []string{"filesystem"},
	"gamelift":               []string{"alias", "build", "fleet"},
	"glue":                   []string{"classifier", "connection", "crawler", "datacatalogencryptionsettings", "database", "devendpoint", "job", "partition", "securityconfiguration", "table", "trigger"},
	"greengrass":             []string{"connectordefinition", "connectordefinitionversion", "coredefinition", "coredefinitionversion", "devicedefinition", "devicedefinitionversion", "functiondefinition", "functiondefinitionversion", "group", "groupversion", "loggerdefinition", "loggerdefinitionversion", "resourcedefinition", "resourcedefinitionversion", "subscriptiondefinition", "subscriptiondefinitionversion"},
	"guardduty":              []string{"detector", "filter", "ipset", "master", "member", "threatintelset"},
	"iam":                    []string{"accesskey", "group", "instanceprofile", "managedpolicy", "policy", "role", "servicelinkedrole", "user", "usertogroupaddition"},
	"inspector":              []string{"assessmenttarget", "assessmenttemplate", "resourcegroup"},
	"iot1click":              []string{"device", "placement", "project"},
	"iot":                    []string{"certificate", "policy", "policyprincipalattachment", "thing", "thingprincipalattachment", "topicrule"},
	"iotanalytics":           []string{"channel", "dataset", "datastore", "pipeline"},
	"kms":                    []string{"alias", "key"},
	"kinesis":                []string{"stream", "streamconsumer"},
	"kinesisanalytics":       []string{"application", "applicationoutput", "applicationreferencedatasource"},
	"kinesisanalyticsv2":     []string{"application", "applicationcloudwatchloggingoption", "applicationoutput", "applicationreferencedatasource"},
	"kinesisfirehose":        []string{"deliverystream"},
	"lambda":                 []string{"alias", "eventsourcemapping", "function", "layerversion", "layerversionpermission", "permission", "version"},
	"logs":                   []string{"destination", "loggroup", "logstream", "metricfilter", "subscriptionfilter"},
	"mediastore":             []string{"container"},
	"neptune":                []string{"dbcluster", "dbclusterparametergroup", "dbinstance", "dbparametergroup", "dbsubnetgroup"},
	"opsworks":               []string{"app", "elasticloadbalancerattachment", "instance", "layer", "stack", "userprofile", "volume"},
	"opsworkscm":             []string{"server"},
	"pinpointemail":          []string{"configurationset", "configurationseteventdestination", "dedicatedippool", "identity"},
	"ram":                    []string{"resourceshare"},
	"rds":                    []string{"dbcluster", "dbclusterparametergroup", "dbinstance", "dbparametergroup", "dbsecuritygroup", "dbsecuritygroupingress", "dbsubnetgroup", "eventsubscription", "optiongroup"},
	"redshift":               []string{"cluster", "clusterparametergroup", "clustersecuritygroup", "clustersecuritygroupingress", "clustersubnetgroup"},
	"robomaker":              []string{"fleet", "robot", "robotapplication", "robotapplicationversion", "simulationapplication", "simulationapplicationversion"},
	"route53":                []string{"healthcheck", "hostedzone", "recordset", "recordsetgroup"},
	"route53resolver":        []string{"resolverendpoint", "resolverrule", "resolverruleassociation"},
	"s3":                     []string{"bucket", "bucketpolicy"},
	"sdb":                    []string{"domain"},
	"ses":                    []string{"configurationset", "configurationseteventdestination", "receiptfilter", "receiptrule", "receiptruleset", "template"},
	"sns":                    []string{"subscription", "topic", "topicpolicy"},
	"sqs":                    []string{"queue", "queuepolicy"},
	"ssm":                    []string{"association", "document", "maintenancewindow", "maintenancewindowtarget", "maintenancewindowtask", "parameter", "patchbaseline", "resourcedatasync"},
	"sagemaker":              []string{"endpoint", "endpointconfig", "model", "notebookinstance", "notebookinstancelifecycleconfig"},
	"secretsmanager":         []string{"resourcepolicy", "rotationschedule", "secret", "secrettargetattachment"},
	"servicecatalog":         []string{"acceptedportfolioshare", "cloudformationproduct", "cloudformationprovisionedproduct", "launchnotificationconstraint", "launchroleconstraint", "launchtemplateconstraint", "portfolio", "portfolioprincipalassociation", "portfolioproductassociation", "portfolioshare", "resourceupdateconstraint", "tagoption", "tagoptionassociation"},
	"servicediscovery":       []string{"httpnamespace", "instance", "privatednsnamespace", "publicdnsnamespace", "service"},
	"stepfunctions":          []string{"activity", "statemachine"},
	"transfer":               []string{"server", "user"},
	"waf":                    []string{"bytematchset", "ipset", "rule", "sizeconstraintset", "sqlinjectionmatchset", "webacl", "xssmatchset"},
	"wafregional":            []string{"bytematchset", "geomatchset", "ipset", "ratebasedrule", "regexpatternset", "rule", "sizeconstraintset", "sqlinjectionmatchset", "webacl", "webaclassociation", "xssmatchset"},
	"workspaces":             []string{"workspace"},
	"ask":                    []string{"skill"},
}

// AllResources returns all possible resources
func AllResources() []string {
	_, resources := AllControllersAndResources()
	return resources
}

// AllControllers returns all possible controllers
func AllControllers() []string {
	ctrls, _ := AllControllersAndResources()
	return ctrls
}

// AllControllersAndResources returns controllers and resources as slices
func AllControllersAndResources() ([]string, []string) {
	controllersResp := []string{}
	resourcesResp := []string{}

	for key, resources := range AvailableControllers {
		controllersResp = append(controllersResp, key)
		for _, resource := range resources {
			resourcesResp = append(resourcesResp, strings.Join([]string{key, resource}, "."))
		}
	}
	return controllersResp, resourcesResp
}
