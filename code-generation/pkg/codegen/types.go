package codegen

// ModelFiles returns all model files in a list format
type ModelFiles struct {
	APIVersion string      `yaml:"apiVersion",json:"apiVersion"`
	Kind       string      `yaml:"kind",json:"kind"`
	Items      []ModelFile `yaml:"items",json:"items"`
}

// ModelFile parses the model files for each service
type ModelFile struct {
	APIVersion string        `yaml:"apiVersion",json:"apiVersion"`
	Kind       string        `yaml:"kind",json:"kind"`
	Metadata   ModelMetadata `yaml:"metadata",json:"metadata"`
	Spec       ModelSpec     `yaml:"spec",json:"spec"`
}

// ModelMetadata stores any specific metadata about the model
type ModelMetadata struct {
	Name string `yaml:"name",json:"name"`
}

// ModelSpec defines the core of how the resource is structured
type ModelSpec struct {
	Kind                string              `yaml:"kind",json:"kind"`
	Type                string              `yaml:"type",json:"type"`
	Queue               bool                `yaml:"queue",json:"queue"`
	UseCloudFormation   bool                `yaml:"useCloudFormation",json:"useCloudFormation"`
	Resource            SpecResource        `yaml:"resource",json:"resource"`
	Body                ResourceBody        `yaml:"body",json:"body"`
	Output              ResourceOutput      `yaml:"output",json:"output"`
	AdditionalResources AdditionalResources `yaml:"additionalResources",json:"additionalResources"`
	Customizations      Customizations      `yaml:"customizations",json:"customizations"`
	IsCustomized        bool
}

// ResourceBody defines the body of the object
type ResourceBody struct {
	Schema SpecSchema `yaml:"schema",json:"schema"`
}

// ResourceOutput defines the body of the object
type ResourceOutput struct {
	Schema SpecSchema `yaml:"schema",json:"schema"`
}

// AdditionalResources defines what is created after the fact
type AdditionalResources struct {
	Services   []Service   `yaml:"services",json:"services"`
	ConfigMaps []ConfigMap `yaml:"configMaps",json:"configMaps"`
	Secrets    []Secret    `yaml:"secrets",json:"secrets"`
}

// SpecResource defines how the CRD is populated
type SpecResource struct {
	Name       string      `yaml:"name",json:"name"`
	Plural     string      `yaml:"plural",json:"plural"`
	Shortnames []Shortname `yaml:"shortNames",json:"shortNames"`
	Scope      string      `yaml:"scope",json:"scope"`
}

// Shortname defines the shortnames the crd will listen to
type Shortname struct {
	Name string `yaml:"name",json"name"`
}

// SpecSchema defines how the object is defined in types.go and the CFT
type SpecSchema struct {
	Type       string           `yaml:"type",json:"type"`
	Properties []SchemaProperty `yaml:"properties",json:"properties"`
	KeyMapping map[string][]SchemaProperty
}

// SchemaProperty is the definition for the full properties
type SchemaProperty struct {
	Key         string           `yaml:"key",json:"key"`
	Type        string           `yaml:"type",json:"type"`
	Description string           `yaml:"description",json:"description"`
	StructKey   string           `yaml:"structKey",json:"structKey"`
	TemplateKey string           `yaml:"templateKey",json:"templateKey"`
	Templatized bool             `yaml:"templatized",json:"templatized"`
	Template    string           `yaml:"template",json:"template"`
	Properties  []SchemaProperty `yaml:"properties",json:"properties"`
}

// Customizations returns the able customizations
type Customizations struct {
	Package string `yaml:"package",json:"package"`
	Add     string `yaml:"add",json:"add"`
	Update  string `yaml:"update",json:"update"`
	Delete  string `yaml:"delete",json:"delete"`
}

// Service defines what you can mutate in a service object
type Service struct {
	Name         string `yaml:"name",json:"name"`
	Type         string `yaml:"type",json:"type"`
	ExternalName string `yaml:"externalName",json:"externalName"`
	Ports        []Port `yaml:"ports",json:"ports"`
}

// Port defines the ServicePorts
type Port struct {
	Port     string `yaml:"port",json:"port"`
	Protocol string `yaml:"protocol",json:"protcol"`
}

// ConfigMap defines what you can mutate in a configmap object
type ConfigMap struct {
	Name string          `yaml:"name",json:"name"`
	Data []ConfigMapData `yaml:"data",json:"data"`
}

// ConfigMapData defines what the data is in a CM.
type ConfigMapData struct {
	Key   string `yaml:"key",json:"key"`
	Value string `yaml:"value",json:"value"`
}

// Secret defines what you can mutate in a secret object
type Secret struct {
}
