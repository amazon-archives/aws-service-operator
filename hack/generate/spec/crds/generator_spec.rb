dir = File.expand_path('../..', __dir__)
$:.unshift(dir) unless $:.include?(dir)

require "crd/generator"

RSpec.describe CRDGenerator, "class" do
  let(:spec) { OpenStruct.new({model_files: []}) }
  context "initializing" do
    it "has region set" do
      crds = CRDGenerator.new("us-east-1", spec)
      expect(crds.region).to eq "us-east-1"
    end
  end
end
