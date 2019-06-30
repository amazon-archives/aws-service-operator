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

class ConfigGenerator
  include Helpers
  attr_reader :region, :version, :models

  def initialize(region, version, models)
    @region = region
    @version = version
    @models = models.dup
  end

  def generate(files=true)
    file_name = ".awsoperator.yaml"
    config = {
      "apiVersion" => "self.awsoperator.io/v1alpha1",
      "kind" => "Config",
      "clusterName" => "aws-service-operator",
      "kubernetes" => {},
      "aws" => {
        "defaultRegion" => "us-west-2",
        "supportedRegions" => ["us-west-2"],
        "accountId" => "XXX",
        "queue" => {
          "name" => "aws-service-operator",
          "region" => "us-west-2"
        }
      },
      "server" => {
        "metrics" => {
          "enabled" => true,
          "endpoint" => "/metrics"
        },
        "log" => {
          "level" => 1
        }
      },
      "resources" => resources
    }

    if files
      puts "storing #{file_name}"
      File.open(file_name, 'w') do |file|
        file.write(config.to_yaml)
      end
    end
  end

  private
  def resources
    list = @models.map do |model|
      [
        model["spec"]["resource"]["apiGroup"]["group"],
        model["spec"]["resource"]["names"]["singular"]
      ].join(".")
    end
    list.push("cloudformation.stack")
  end
end