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

class ExamplesGenerator
  include Helpers
  attr_reader :spec, :examples_dir

  def initialize(spec, examples_dir="examples")
    @spec = spec
    @examples_dir = examples_dir
  end

  def generate
    puts "Removing example directory"
    FileUtils.rm_rf("#{@examples_dir}/")

    @spec.json["ResourceTypes"].each do |resource_name, resource|
      name_split = resource_name.split("::")
      group_name = name_split[1]
      type_name = resource_name = name_split.last
      cr = {
        "apiVersion" => "#{group_name.downcase}.awsoperator.io/v1alpha1",
        "kind" => type_name.camelize,
        "metadata" => {
          "name" => "example-#{group_name.downcase}-#{type_name.downcase}",
          "namespace" => "default",
        },
        "spec" => generate_properties(resource, resource_name, @spec.json, resource_name)
      }
      cr["spec"].delete("name")

      puts "Making example directory and subfolder #{group_name.downcase}"
      FileUtils.mkdir_p("#{@examples_dir}/#{group_name.downcase}")

      puts "Storing #{@examples_dir}/#{group_name.downcase}/#{type_name.downcase}.yaml"
      File.open("#{@examples_dir}/#{group_name.downcase}/#{type_name.downcase}.yaml", 'w') { |file| file.write(cr.to_yaml()) }
    end
  end

  private
  def reference_object(resource_name, type)
    obj = {
      "arn" => "String",
      "#{resource_name.singularize.camelize(:lower)}Ref" => {
        "name" => "String",
        "namespace" => "String"
      }
    }
    if type == "List"
      return [obj]
    else
      return obj
    end
  end

  def generate_property(property_name, property, resource, resource_name, json, root_name)
    if property["PrimitiveType"]
      return property["PrimitiveType"]
    elsif property["Type"] == "List" && property["PrimitiveItemType"]
      return [property["PrimitiveItemType"]]
    elsif property["Type"] == "List" && property["ItemType"]
      if property["ItemType"] != "Tag"
        return [generate_properties(json["PropertyTypes"]["#{root_name}.#{property["ItemType"]}"], property_name, json, root_name)]
      else
        return [{"key" => "String", "value" => "String"}]
      end
    elsif property["Type"] == "Map"
      return {}
    elsif property["Type"]
      return generate_properties(json["PropertyTypes"]["#{root_name}.#{property["Type"]}"], property_name, json, root_name)
    elsif property["Type"]
      return property["PrimitiveType"]
    end
  end

  def generate_properties(resource, resource_name, json, root_name)
    ret = {}

    begin
      resource["Properties"].each do |property_name, property|
        if !resource_reference?(property_name)
          if property_name != resource_name
            ret[sanatize_property_name(property_name, resource_name, root_name)] = generate_property(property_name, property, resource, resource_name, json, root_name)
          else
            ret[sanatize_property_name(property_name, resource_name, root_name)] = {}
          end
        else
          name = property_name.gsub(/Arn$/, '').camelize(:lower)
          ret[name] = reference_object(name, property["Type"])
        end
      end
    rescue
      ret[sanatize_property_name(resource_name, resource_name, root_name)] = resource
    end


    return ret
  end

end