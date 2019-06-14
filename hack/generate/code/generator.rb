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

require "open-uri"
require "git"

class CodeGenerator
  attr_accessor :json
  attr_reader :region, :repo_url, :name, :repo_dir

  def initialize(region, repo_url="git@github.com:kubernetes/code-generator.git", name="code-generator", repo_dir="hack/generate")
    @region = region
    @repo_url =  repo_url
    @name = name
    @repo_dir = repo_dir
  end

  def clone
    if File.directory?("#{repo_dir}/#{name}")
      g = Git.open(File.join(repo_dir, name))
      g.pull('origin', 'release-1.13')
    else
      g = Git.clone(@repo_url, @name, {path: repo_dir})
      g.checkout('release-1.13')
      g.pull('origin', 'release-1.13')
    end
  end

  def clean_repo
    FileUtils.rm_rf(File.join(repo_dir, name)) unless ENV["FAST_MODE"]
  end

  def generate_groups(commands, groups_list)
    system("./hack/generate/code-generator/generate-groups.sh #{commands} awsoperator.io/pkg/generated awsoperator.io/pkg/apis \"#{groups_list}\" --output-base #{ENV["GOPATH"]}src --go-header-file hack/generate/templates/boilerplate.go.txt")
  end
end