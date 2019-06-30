# Copyright 2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
#     http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.

require "helpers/helpers"

class CRDGenerator
  include Helpers
  attr_reader :region, :models, :config_dir, :crds

  def initialize(region, models, config_dir="configs")
    @region = region
    @models = models.dup
    @config_dir = config_dir
    @crds = []
  end

  def generate(files=true)
    if files
      puts "Making CRDs config directory"
      config_directory = File.join("configs", "crds")
      FileUtils.rm_rf(config_directory)
      FileUtils.mkdir_p(config_directory)
    end

    @models.each do |model|
      crd = {
        "apiVersion" => "apiextensions.k8s.io/v1beta1",
        "kind" => "CustomResourceDefinition",
        "metadata" => {
          "name" => model["spec"]["resource"]["names"]["plural"]+"."+group_joined(model, "."),
        },
        "spec" => {
          "group" => group_joined(model, "."),
          "versions" => versions(model),
          "scope" => model_scope(model),
          "names" => {
            "plural" => model["spec"]["resource"]["names"]["plural"],
            "singular" => model["spec"]["resource"]["names"]["singular"],
            "kind" => model["spec"]["resource"]["kind"],
            "shortNames" => model["spec"]["resource"]["names"]["shortNames"],
            "categories" => [
              "aws",
              model["spec"]["resource"]["apiGroup"]["group"]
            ]
          },
          "additionalPrinterColumns" => [
            {
              "name" => "Status",
              "type" => "string",
              "description" => "Status for the AWS CloudFormation stack.",
              "JSONPath" => ".status.reason"
            },
            {
              "name" => "Message",
              "type" => "string",
              "description" => "Message accompanying it's current stack status.",
              "JSONPath" => ".status.message",
              "priority" => 1
            },
            {
              "name" => "Created At",
              "type" => "date",
              "description" => "When the resource was created",
              "JSONPath" => ".metadata.creationTimestamp"
            }
          ],
          "subresources" => {
            "status" => {}
          },
          "validation" => {
            "openAPIV3Schema" => {
              "properties" => {
                "spec" => open_api_v3_schema(model["spec"]["openAPIV3Schema"]["properties"]["spec"])
              }
            }
          }
        }
      }

      @crds << crd

      if files
        model_group = model["spec"]["resource"]["apiGroup"]["group"]
        directory_path = File.join(config_directory, model_group)
        FileUtils.mkdir_p(directory_path)

        model_group = model["spec"]["resource"]["names"]["singular"]
        file_name = File.join(directory_path, "#{model_group}.yaml")
        puts "storing #{file_name}"
        File.open(file_name, 'w') do |file|
          file.write(crd.to_yaml)
        end
      end
    end

    FileUtils.copy_entry("hack/generate/templates/crds/cloudformation/", File.join(config_directory, "cloudformation"))
  end

  private
  def open_api_v3_properties(spec)
    props = {}
    spec.each do |key, value|
      # binding.pry if value["description"] == "http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-amazonmq-broker.html#cfn-amazonmq-broker-securitygroups"
      props[key] = open_api_v3_schema(value)
    end
    props.delete_if {|k,v| v.nil? }
    props
  end

  def open_api_v3_schema(schema)
    prop = schema.dup
    return nil if prop["customAttributes"] && prop["customAttributes"]["resourceName"]
    prop.delete("customAttributes")

    prop["properties"] = open_api_v3_properties(prop["properties"]) if prop["properties"]
    prop["items"] = open_api_v3_properties(prop["items"]) if prop["items"]
    prop
  end

  def model_scope(model)
    model["spec"]["resource"]["scope"]
  end

  def versions(model)
    model["spec"]["resource"]["apiGroup"]["versions"].map do |version|
      {
        "name" => version["name"],
        "served" => version["served"],
        "storage" => version["storage"],
      }
    end
  end

  def group_joined(model, joiner)
    [model["spec"]["resource"]["apiGroup"]["group"],
     model["spec"]["resource"]["apiGroup"]["domain"]].join(joiner)
  end
end