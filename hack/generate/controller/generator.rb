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

require 'erb'
require 'open-uri'
require "helpers/helpers"
require "api/model"

class ControllerGenerator
  include Helpers
  attr_reader :region, :version, :models, :controller_dir

  def initialize(region, version, models, controller_dir="pkg/generated/controllers")
    @region = region
    @version = version
    @models = models.dup
    @controller_dir = controller_dir
  end

  def generate(files=true)
    controller_directory = File.join(@controller_dir)
    if files
      puts "Making the controllers directory"
      FileUtils.rm_rf(controller_directory)
      FileUtils.mkdir_p(controller_directory)
    end
    FileUtils.copy_entry("hack/generate/templates/controllers/", File.join(controller_dir))


    # Generate controllers for other resources
  end

  private
  def template(file, bind)
    content = open(file)
    erb = ERB.new(content.read)
    erb.result(bind)
  end
end