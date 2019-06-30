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

require "hana"
require "json"
require "open-uri"

class PatchGenerator
  attr_reader :spec, :patch_dir, :patched_spec

  def initialize(spec, patch_dir="hack/generate/customizations/specs")
    @spec = spec
    @patch_dir = patch_dir
  end

  # Any time we need to patch new type add the methods here
  def generate(files=true)
    generate_scope_patch
  end

  def apply
    @patched_spec = @spec
    Dir["#{patch_dir}/*.json"].each do |file|
      content = open(file)
      json = JSON.parse(content.read)
      patch = Hana::Patch.new(json)
      @patched_spec.json = patch.apply(@patched_spec.json)
    end
    @patched_spec
  end

  private
  def generate_scope_patch
    patches = @spec.json["ResourceTypes"].keys.map do |name|
      if name =~ /AWS::IAM::/
      {
          "op": "add",
          "path": "/ResourceTypes/#{name}/Scope",
          "value": {
            # "ScopeType": "Cluster"
            "ScopeType": "Namespaced"
          }
        }
      else
        {
            "op": "add",
            "path": "/ResourceTypes/#{name}/Scope",
            "value": {
              "ScopeType": "Namespaced"
            }
          }
      end
    end

    File.open(File.join(patch_dir, "02_scope_values.json"), "w") do |file|
      file.write(JSON.pretty_generate(patches))
    end
  end

end