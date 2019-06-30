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

require "active_support/core_ext/string"
require "helpers/helpers"
require "cfnspec/cfnspec"
require "patch/generator"
require "model_file/shame_list"
require "model_file/shared_resources"


class ModelFileGenerator
  include Helpers
  include ShameList
  include SharedResources
  attr_accessor :model_files, :spec, :sdk_models
  attr_reader :region, :models_dir, :domain, :version, :patcher

  TYPES = %w(integer, number, string, boolean, object)
  FORMATS = %w(int32, int64, float, double, byte, binary, date, date-time, password)

  def initialize(region, domain="awsoperator.io", version="v1alpha1", models_dir="models")
    @region = region
    @spec = initialize_spec
    @models_dir = models_dir
    @domain = domain
    @version = version
    @model_files = []

    generate_patches
    apply_patches
  end

  def generate(files=true)
    if files
      puts "removing models directory"
      dir = @models_dir + "/"
      FileUtils.rm_rf(dir)
      FileUtils.mkdir_p(dir)
    end

    @spec.json["ResourceTypes"].each do |resource_name, resource|
      group_name, type_name = get_group_and_type_name(resource_name)
      
      # We handle CFN resources differently
      next if group_name.downcase == "cloudformation"

      model_file = {
        "apiVersion" => "models.awsoperator.io/v1alpha1",
        "kind" => "ModelFile",
        "metadata" => {
          "name" => [group_name.downcase, type_name.downcase, "resource"].join("-"),
          "creationTimestamp" => Time.now.iso8601,
          "resourceVersion" => @spec.json["ResourceSpecificationVersion"],
        },
        "spec" => {
          "resource" => {
            "kind" => kind_name(type_name, group_name),
            "scope" => @spec.json["ScopeTypes"][resource["Scope"]["ScopeType"]]["Value"],
            "apiGroup" => {
              "apiVersion" => "#{group_name.downcase}.#{@domain}/#{@version}",
              "group" => group_name.downcase,
              "region" => @region,
              "domain" => @domain,
              # TODO: Consider how to better version these
              "versions" => [
                {
                  "name" => @version,
                  "served" => true,
                  "storage" => true
                }
              ]
            },
            "names" => {
              "singular" => type_name.downcase,
              "plural" => type_name.downcase.pluralize,
              "shortNames" => []
            }
          },
          "openAPIV3Schema" => {
            "properties" => {
              "spec" => generate_property(resource_name, "spec", resource)
            }
          }
        }
      }

      # Add to @model_files so we can use this object directly instead of
      # parsing again
      @model_files << model_file

      if files
        directory_path = File.join(@models_dir, group_name.downcase)
        FileUtils.mkdir_p(directory_path)

        file_name = File.join(directory_path, "#{type_name.downcase}.yaml")
        puts "storing #{file_name}"
        File.open(file_name, 'w') do |file|
          file.write(model_file.to_yaml)
        end
      end
    end
  end

  private
  def kind_name(type_name, group_name)
    name = type_name.camelize
    name = [group_name, name].map(&:camelize).join if name == "Resource"
    name
  end

  def primitive_type_and_format(prim_type)
    type = prim_type.downcase
    resp = case type
    when "string"
      ["string", ""]
    when "boolean"
      ["boolean", ""]
    when "integer"
      ["integer", "int32"]
    when "double"
      ["integer", "double"]
    when "long"
      ["integer", "int64"]
    end
    resp
  end

  def type_and_format(prim_type, type)
    if prim_type
      return primitive_type_and_format(prim_type)
    end

    if type == "List"
      return ["array", ""]
    end

    if type == "Map"
      return ["json", ""]
    end

    return ["object", ""]
  end

  def id?(key)
    key.length > 2 && key =~ /Id|Ids$/
  end

  def arn?(key)
    key.length > 3 && key =~ /Arn|Arns$/
  end

  def security_group?(key)
    key =~ /^securitygroup/i
  end

  def id_arn_security_group?(key)
    id?(key) || arn?(key) || security_group?(key)
  end

  def id_arn_sg?(key, resource)
    return if resource["PrimitiveType"] != "String"

  end

  def plural?(key)
    key.pluralize == key && key.singularize != key
  end

  def resource_name?(resource_name, resource_key, resource)
    group_type_name = get_group_and_type_name(resource_name)
    true if resource["UpdateType"] == "Immutable" && !!(generate_key(resource_key) =~ /^#{generate_key(group_type_name.last)}Name$/)
  end

  def generate_key(key)
    resp = key
    resp = key == "type" ? "contentType" : key
    if id_arn_security_group?(key)
      resp = key.gsub(/arn|s$/i, "") if arn?(key)
      resp = key.gsub(/id|s$/i, "") if id?(key)
      resp = resp+"Ref"
    end
    resp = key if resp == ""

    resp.camelize(:lower)
  end

  def key_type(key)
    if id?(key)
      "id"
    elsif arn?(key)
      "arn"
    elsif security_group?(key)
      "securityGroup"
    end
  end

  def reference_object(resource, resource_key)
    key_type = key_type(resource_key)
    prop = {}
    prop["description"] = resource["Documentation"]
    prop["customAttributes"] = reference_custom_attrs(resource, resource_key)
    if plural?(resource_key)
      prop["type"] = "array"
      prop["items"] = reference_property(key_type.pluralize, resource_key)
    else
      prop["type"] = "object"
      prop["properties"] = reference_property(key_type, resource_key)
    end
    prop
  end

  def reference_custom_attrs(resource, resource_key)
    prop = {}
    prop["template"] = resource_key
    prop["immutable"] = if resource["UpdateType"] == "Immutable"
      true
    elsif resource["UpdateType"] == "Mutable"
      false
    end
    prop
  end

  def reference_property(key_type, resource_key)
    {
      "#{key_type}" => {
        "type" => "string",
        "description" => "raw #{key_type} for the #{resource_key}"
      },
      "#{generate_key(resource_key)}Ref" => {
        "type" => "object",
        "description" => "#{resource_key} reference using other CRDs",
        "properties" => {
          "name" => {
            "type" => "string",
            "description" => "#{resource_key} name reference for other CRD"
          },
          "namespace" => {
            "type" => "string",
            "description" => "#{resource_key} namespace reference for other CRD"
          }
        }
      }
    }
  end

  def array_items(resource_name, resource)
    proptype, propformat = type_and_format(resource["PrimitiveItemType"], resource["Type"])

    prop = {}
    prop["type"] = proptype

    if proptype == "array"
      property_type = @spec.json["PropertyTypes"][resource_name + "." + resource["ItemType"]]

      if property_type.nil? && resource["ItemType"] == "Tag"
        property_type = resource_tag_property_type
      end

      properties = !RECURSIVE_SHAME_LIST.include?(resource_name + "." + resource["ItemType"])
      prop["properties"] = generate_properties(resource_name, property_type, properties)
      prop["type"] = "object"
    end

    prop["format"] = propformat

    prop.delete_if {|key, value| value.to_s.empty? }
    return prop
  end

  def generate_properties(resource_name, property_type, properties=true, nested=[])
    props = {}
    if properties && !property_type["Properties"].nil?
      property_type["Properties"].each do |k, v|
        props[generate_key(k)] = generate_property(resource_name, k, v, nested)
      end
    end
    props
  end

  def generate_property(resource_name, resource_key, resource, nested=[])
    prop = {}
    nested << resource_key unless resource_key == "spec"

    if id_arn_sg?(resource_key, resource)
      prop = reference_object(resource, resource_key)
    else
      proptype, propformat = type_and_format(resource["PrimitiveType"], resource["Type"])

      prop["description"] = resource["Documentation"]
      prop["type"] = proptype
      prop["format"] = propformat
      prop["customAttributes"] = {}
      prop["customAttributes"]["template"] = resource_key unless proptype == "object"

      unless resource["Properties"] && !resource_name?(resource_name, resource_key, resource)
        prop["customAttributes"]["resourceName"] = resource_name?(resource_name, resource_key, resource)
      end

      # technically there is three states, and rather than case statement only
      # account for 2 values
      prop["customAttributes"]["immutable"] = if resource["UpdateType"] == "Immutable"
        true
      elsif resource["UpdateType"] == "Mutable"
      false
      end

      prop["customAttributes"].delete_if {|key, value| value.to_s.empty? }

      if resource["Type"] && resource["Type"] != "List" && resource["Type"] != "Map"
          property_type = @spec.json["PropertyTypes"][resource_name + "." + resource["Type"]]

          unless property_type["Properties"].nil?
            prop["required"] = property_type["Properties"].map {|k, v| k if v["Required"]}.compact.map {|v| v.camelize(:lower) }
            prop["properties"] = generate_properties(resource_name, property_type, true, nested)
          end
      elsif resource["Properties"]
        prop["required"] = resource["Properties"].map do |k, v|
          k if v["Required"] if k.camelize(:lower) != (resource_name.split("::").last+"Name").camelize(:lower)
        end.compact.map do |v|
          v.camelize(:lower)
        end.delete_if(&:nil?)
        prop["properties"] = generate_properties(resource_name, resource, true, nested)
      end

      if proptype == "array"
        prop["items"] = array_items(resource_name, resource)
      end

      resource_value = resource["Value"]
      unless resource_value.nil?
        value_type = @spec.json["ValueTypes"][resource_value["ValueType"]]
        unless value_type.nil? && value_type["AllowedValues"].nil?
          prop["enum"] = value_type["AllowedValues"]
        end
      end

      # binding.pry if resource["Documentation"] == "http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicediscovery-service-dnsconfig.html#cfn-servicediscovery-service-dnsconfig-dnsrecords"

      prop.delete_if {|key, value| value.to_s.empty? }
      prop
    end
  end

  def initialize_spec
    spec = CFNSpec.new(@region)
    spec.clean_repo
    spec.clone
    spec.parse
    spec
  end

  def generate_patches
    @patcher = PatchGenerator.new(@spec)
    patcher.generate
  end

  def apply_patches
    @spec = @patcher.apply
  end
end