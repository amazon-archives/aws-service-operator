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
require "code/generator"

class APIGenerator
  include Helpers
  attr_reader :region, :version, :models, :config_dir, :apis

  def initialize(region, version, models, apis_dir="pkg/apis")
    @region = region
    @version = version
    @models = models.dup
    @apis_dir = apis_dir
    @apis = []
  end

  def generate(files=true)
    apis_directory = File.join(@apis_dir)
    if files
      puts "Making the apis directory"
      FileUtils.rm_rf(apis_directory)
      FileUtils.mkdir_p(apis_directory)
    end
    domain = @models.first["spec"]["resource"]["apiGroup"]["domain"]

    groups = {}
    @models.each do |model|
      group = model["spec"]["resource"]["apiGroup"]["group"]
      if !groups[group]
        groups[group] = []
      end

      groups[group] << Model.new(model, domain, group, @version)
    end

    groups.each do |group, model_list|
      ["doc.go", "register.go", "types.go"].each do |file|
        directory_path = File.join(apis_directory, group, @version)
        model_file = Model.new(model_list, domain, group, @version)

        if files
          FileUtils.mkdir_p(directory_path)

          content = template("hack/generate/templates/api_#{file}.erb", model_file.get_binding)

          file_name = File.join(directory_path, file)
          puts "storing #{file_name}"
          File.open(file_name, 'w') { |f| f.write(content) }
          system("goimports -w #{file_name}")
        end
      end
    end

    # Create apis/doc.go
    model_file = Model.new({}, domain, "", @version)
    content = template("hack/generate/templates/root_register.go.erb", model_file.get_binding)

    file_name = File.join(apis_directory, "register.go")
    puts "storing #{file_name}"
    File.open(file_name, 'w') { |f| f.write(content) }
    system("goimports -w #{file_name}")

    # Copy apis/meta
    FileUtils.copy_entry("hack/generate/templates/meta/", File.join(apis_directory, "meta"))
    FileUtils.copy_entry("hack/generate/templates/self/", File.join(apis_directory, "self"))
    FileUtils.copy_entry("hack/generate/templates/cloudformation/", File.join(apis_directory, "cloudformation"))

    content = template("hack/generate/templates/self/v1alpha1/resources.go", AllModels.new(groups).get_binding)
    file_name = File.join(apis_directory, "/self/v1alpha1/", "resources.go")
    puts "storing #{file_name}"
    File.open(file_name, 'w') { |f| f.write(content) }
    system("goimports -w #{file_name}")

    codegen = CodeGenerator.new(@region)
    codegen.clone

    # run generate_groups via shelling out.
    groups_list = groups.map {|key, value| key + ":" + @version }
    groups_list.push("cloudformation:#{@version}")
    codegen.generate_groups("deepcopy,client,informer,lister", groups_list.join(" "))
    codegen.generate_groups("deepcopy", ["meta:#{@version}", "self:#{@version}"].join(" "))
  end

  private
  def template(file, bind)
    content = open(file)
    erb = ERB.new(content.read)
    erb.result(bind)
  end

end