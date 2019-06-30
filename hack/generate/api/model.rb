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

require 'ostruct'

class AllModels
  attr_accessor :obj
  def initialize(obj)
    @obj = obj
  end

  def get_binding
    binding
  end
  
  def available_controllers
    # // "apigateway": []string{"account", "apikey"},
    controllers = @obj.map do |key, value|
      string_slice = value.map {|v| "\"#{v.obj.spec.resource.kind.downcase}\"" }
      "\"#{key}\":[]string{#{string_slice.join(", ")}},"
    end
    controllers.join("\n")
  end
end

class Model
  attr_accessor :raw, :obj, :domain, :group, :version

  def initialize(obj, domain, group, version)
    @obj = obj.is_a?(Array) ? obj : to_recursive_ostruct(obj)
    @domain = domain
    @group = group
    @version = version
  end

  def get_binding
    binding
  end

  def property_type_or_nested(prop, nested)
    resp = case [prop.type.to_s, prop.format.to_s]
    when ["string", ""]
      "string"
    when ["boolean", ""]
      "bool"
    when ["integer", "int32"]
      "int32"
    when ["integer", "double"]
      "float64"
    when ["integer", "int64"]
      "int64"
    when ["json", ""]
      "map[string]string"
    when ["array", ""]
      if ["string", "bool", "int32", "int64", "float64"].include?(prop.items.type)
        "[]#{prop.items.type}"
      else
        "[]#{nested.singularize}"
      end
    else
      nested
    end
    resp
  end

  def primitive_type?(key, prop)
    ["string", "bool", "int32", "int64", "float64", "map[string]string", "[]string", "[]int"].include?(property_type(key, prop, ""))
  end

  def property_type(key, prop, additional)
    resp = case [prop.type.to_s, prop.format.to_s]
    when ["string", ""]
      "string"
    when ["boolean", ""]
      "bool"
    when ["integer", "int32"]
      "int32"
    when ["integer", "double"]
      "float64"
    when ["integer", "int64"]
      "int64"
    when ["json", ""]
      "map[string]string"
    when ["array", ""]
      if ["string", "bool", "int32", "int64", "float64"].include?(prop.items.type)
        "[]#{prop.items.type}"
      else
        "[]#{[self.spec.resource.kind, "Spec", key.to_s.singularize.camelize].join}"
      end
    else
      [self.spec.resource.kind, "Spec", key.to_s.camelize].join
    end
    resp
  end

  def new_object(name, description, properties)
    OpenStruct.new({
      name: name,
      description: description,
      properties: properties
    })
  end

  def new_property(key, description, inline, type, fmt, required, template_tag)
    OpenStruct.new({
      key: key.to_s,
      description: description.to_s,
      inline: inline,
      type: type.to_s,
      format: fmt.to_s,
      required: required,
      template_tag: template_tag
    })
  end

  def property_definition(property)
    lines = []
    if property.inline
      lines << [property.type, "`json:\",inline\"`"].join
    else
      lines << "// \"+optional\"" if property.required
      lines << "// #{property.key.camelize} #{property.description}"
      lines << "#{property.key.camelize} #{property.type} `json:\"#{property.key},omitempty\"#{property.template_tag}`"
    end
    lines
  end

  def type_definition(item)
    lines = []
    name = item.name.gsub("[]", "")
    lines << "// #{name} #{item.description}"
    lines << "type #{name} struct {"
    item.properties.to_a.each do |property|
      lines << property_definition(property)
      lines << ""
    end
    lines << "}"
    lines
  end

  def loop_through_properties(lines, kind, properties)
    properties.each do |key, property|
      next if primitive_type?(key, property)

      object_name = property_type(key, property, "")
      props = property.properties.to_h.map { |k, prop|
        new_property(k, prop.description, false, property_type_or_nested(prop, [object_name, k.to_s.camelize].join), prop.format, false, template_struct_tag(prop))
      }.compact
      object = new_object(object_name, "is the definition for a #{kind} resource", props)

      lines << type_definition(object)
    end
  end

  def compact_nested_properties(nested, props, properties)
    properties.each do |key, property|
      root_key = property.type == "array" ? key.to_s.singularize : key.to_s
      new_key = [nested.to_s.camelize, root_key.camelize].join.to_sym
      props[new_key] = property
      unless property.properties.nil?
        props = compact_nested_properties(new_key, props, property.properties.to_h)
      end
      if !property.items.nil? && property.items.type == "object"
        props = compact_nested_properties(new_key, props, property.items.properties.to_h)
      end
    end
    props
  end

  def types_definitions
    lines = []
    properties = {}

    # Base Spec
    cfn = [new_property("", "embeds Cloudformation specific details", true, "metav1alpha1.CloudFormationMeta", "string", true, "")]
    props = @obj.spec.openAPIV3Schema.properties.spec.properties.to_h.map { |key, property|
      new_property(key, property.description, false, property_type(key, property, ""), property.format, false, template_struct_tag(property))
    }.compact
    object = new_object("#{@obj.spec.resource.kind}Spec", "is the spec for the #{@obj.spec.resource.kind} resource", cfn+props)

    lines << type_definition(object)

    props = compact_nested_properties("", properties, @obj.spec.openAPIV3Schema.properties.spec.properties.to_h)
    properties.delete_if {|k, v| properties.has_key?(k.to_s.singularize.to_sym) && k.to_s.singularize != k.to_s }

    # All attributes
    loop_through_properties(lines, @obj.spec.resource.kind, properties)

    lines.join("\n")
  end

  def template_struct_tag(prop)
    " cloudformation:\"#{prop.customAttributes.template},Parameter\"" if prop.try(:customAttributes).try(:template)
  end

  def method_missing(m, *args, &block)
    if ["get_binding", "property_type", "template_key", "raw", "types_definitions", "primitive_type?", "available_controllers"].include?(m.to_s)
      super(m, *args, &block)
    else
      @obj.send(m, *args, &block)
    end
  end

  private
  def to_recursive_ostruct(hash)
    OpenStruct.new(hash.each_with_object({}) do |(key, val), memo|
      memo[key] = val.is_a?(Hash) ? to_recursive_ostruct(val) : val
    end)
  end

end