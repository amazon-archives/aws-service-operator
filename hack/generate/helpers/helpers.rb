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

module Helpers
  STRIP_LIST = /Amazon|AWS/i
  def resource_name_property(name, root_name)
    return root_name.split("::").last
  end

  def sanatize_property_name(name, resource_name, root_name)
    root_resource_name = resource_name_property(name, root_name)
    resp = name.gsub(/Configurations|Configuration|#{resource_name}|#{resource_name}s|#{root_resource_name}|#{root_resource_name}s/i, '')
    resp = resp == "" ? name : resp
    return resp.camelize(:lower)
  end

  def resource_reference?(name)
    return (/arn$/i =~ name) != nil
  end

  def reference_type(name)
    return name.match(/([a-zA-Z]*)arn$/i).captures
  end

  def strip_brand(str)
    return str.gsub(STRIP_LIST, "")
  end

  def get_group_and_type_name(resource_name)
    name_split = resource_name.split("::")
    group_name = strip_brand(name_split[1])
    type_name = strip_brand(name_split.last)
    return [group_name, type_name]
  end

end