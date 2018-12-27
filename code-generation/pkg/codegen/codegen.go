package codegen

//go:generate go-bindata -pkg $GOPACKAGE -prefix assets/ -o templates.go assets/

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
)

// TypeJSON returns the full JSON response
func (s ModelSpec) TypeJSON() string {
	if s.Type == "Data" {
		return "data"
	}
	return "spec"
}

// PluralName returns the plural version of the name
func (s ModelSpec) PluralName() string {
	pluralExceptions := map[string]string{
		"Endpoints": "Endpoints",
	}
	pluralNamer := namer.NewPublicPluralNamer(pluralExceptions)

	spec := &types.Type{
		Name: types.Name{
			Name: s.Kind,
		},
	}
	return pluralNamer.Name(spec)
}

// TypeOfObject returns the Type of Object it is
func (s SchemaProperty) TypeOfObject(kind string) string {
	if s.Type == "object" {
		return kind + s.StructKey
	}
	return s.Type
}

// ServiceCount is the amount of services
func (a AdditionalResources) ServiceCount() int {
	return len(a.Services)
}

// ConfigMapCount is the amount of configmaps
func (a AdditionalResources) ConfigMapCount() int {
	return len(a.ConfigMaps)
}

// SecretCount is the amount of secrets
func (a AdditionalResources) SecretCount() int {
	return len(a.Secrets)
}

// NameToLowerCamel will lowercase the name for variables in golang
func (s Service) NameToLowerCamel() string {
	return strcase.ToLowerCamel(s.Name)
}

// Codegen hold the base contructs
type Codegen struct {
	ModelPath string
	RootPath  string
}

// New will create a new Codegen pointer with a config
func New(modelPath string, rootPath string) *Codegen {
	return &Codegen{
		ModelPath: modelPath,
		RootPath:  rootPath,
	}
}

// Run will run the loop over the model files
func (c *Codegen) Run() error {
	modelPath := c.ModelPath
	rootPath := c.RootPath
	log.Infof("Model Path: %+v", modelPath)

	modelFiles, err := ioutil.ReadDir(modelPath)
	if err != nil {
		log.Fatal("Cannot find any model files")
	}

	models := ModelFiles{}

	for _, f := range modelFiles {
		if strings.Contains(f.Name(), ".yaml") {
			parsedModel := ModelFile{}
			log.Info(f.Name())
			modelFileContents, err := ioutil.ReadFile(modelPath + f.Name())

			if err != nil {
				log.Fatalf("Errored reading file %s with error %+v", f.Name(), err)
			}
			err = yaml.Unmarshal(modelFileContents, &parsedModel)
			if err != nil {
				log.Fatalf("Error parsing file %s errored with %+v", f.Name(), err)
			}
			models.Items = append(models.Items, parsedModel)

			operatorPath := rootPath + "pkg/operators/" + parsedModel.Spec.Resource.Name
			apiPath := rootPath + "pkg/apis/service-operator.aws/v1alpha1"

			createDirIfNotExist(operatorPath)

			parsedModel.Spec.Body.Schema.KeyMapping = remapSchema(parsedModel.Spec.Body.Schema.Properties, parsedModel.Spec.Type)
			// parsedModel.Spec.Output.Schema.KeyMapping = remapSchema(parsedModel.Spec.Output.Schema.Properties)
			parsedModel.Spec.IsCustomized = (Customizations{}) == parsedModel.Spec.Customizations

			if (Customizations{}) == parsedModel.Spec.Customizations {
				createFile(rootPath, "cft.go", "cft.go", operatorPath+"/", parsedModel)
			}
			createFile(rootPath, "operator.go", "operator.go", operatorPath+"/", parsedModel)

			createFile(rootPath, parsedModel.Spec.Resource.Name+".go", "types.go", apiPath+"/", parsedModel)

		}
	}

	helpersPath := rootPath + "pkg/helpers"
	configPath := rootPath + "configs"
	operatorsPath := rootPath + "pkg/operators/base"

	createFile(rootPath, "template_functions.go", "template_functions.go", helpersPath+"/", models)
	createFile(rootPath, "aws-service-operator.yaml", "aws-service-operator.yaml", configPath+"/", models)
	createFile(rootPath, "base.go", "base.go", operatorsPath+"/", models)

	return nil
}

func remapSchema(properties []SchemaProperty, typeOf string) map[string][]SchemaProperty {
	remap := make(map[string][]SchemaProperty)
	for _, p := range properties {
		if len(p.Properties) > 0 {
			for _, np := range p.Properties {
				if _, test := remap[p.StructKey]; test {
					remap[p.StructKey] = append(remap[p.StructKey], np)
				} else {
					remap[p.StructKey] = []SchemaProperty{np}
				}
			}
		}

		if _, test := remap[typeOf]; test {
			remap[typeOf] = append(remap[typeOf], p)
		} else {
			remap[typeOf] = []SchemaProperty{p}
		}
	}

	return remap
}

func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatalf("Error creating directories: %+v", err)
		}
	}
}

func createFile(rootPath string, fileName string, templateName string, path string, data interface{}) {
	t := template.New(rootPath + "pkg/codegen/assets/" + templateName + ".templ")
	asset, err := Asset(templateName + ".templ")
	if err != nil {
		log.Fatalf("Error fetching template file: %+v", err)
		return
	}
	t, err = t.Parse(string(asset))
	if err != nil {
		log.Fatalf("Error parsing templates with error: %+v", err)
		return
	}

	bf := bytes.NewBuffer([]byte{})
	err = t.Execute(bf, data)
	if err != nil {
		log.Print("execute: ", err)
		return
	}

	formatted := bf.Bytes()
	if templateName[len(templateName)-3:] == ".go" {
		formatted, err = format.Source(bf.Bytes())
		if err != nil {
			log.Fatalf("Error formatting resolved template %s: %+v", templateName, err)
		}
	}

	f, err := os.Create(path + fileName)
	if err != nil {
		log.Fatalf("Error creating file with error: %+v", err)
		return
	}

	if _, err := f.Write(formatted); err != nil {
		log.Fatalf("Error writing file: %+v", err)
	}

	f.Close()
	return
}
