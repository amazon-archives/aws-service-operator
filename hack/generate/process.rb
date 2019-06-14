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

require 'fileutils'
require "open-uri"
require "json"
require "yaml"
require "active_support/core_ext/string"

dir = File.expand_path('.', __dir__)
$:.unshift(dir) unless $:.include?(dir)

require "model_file/generator"
require "examples/generator"
require "crd/generator"
require "api/generator"
require "config/generator"
require "controller/generator"

files = ENV["FILES"] != "false"
version = ENV["VERSION"] || "v1alpha1"

# This should be set as an ENV for testing
regions = ["us-east-1"]

regions.each do |region|
  model_file = ModelFileGenerator.new(region)
  model_file.generate(files)

  # Create CRDs
  crds = CRDGenerator.new(region, model_file.model_files)
  crds.generate(files)

  # Create APIs
  apis = APIGenerator.new(region, version, model_file.model_files)
  apis.generate(files)

  # Create Config
  config = ConfigGenerator.new(region, version, model_file.model_files)
  config.generate(files)

  # Generate Operators
  controllers = ControllerGenerator.new(region, version, model_file.model_files)
  controllers.generate(files)

  # # Create Examples
  # example = ExamplesGenerator.new(model_file.model_files)
  # example.generate
end

