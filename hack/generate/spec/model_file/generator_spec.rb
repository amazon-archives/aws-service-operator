dir = File.expand_path('../..', __dir__)
$:.unshift(dir) unless $:.include?(dir)

require "model_files/generator"

RSpec.describe ModelFilesGenerator, "class" do
  context "with region" do
    it "generates new modelfiles" do
      spec = ModelFilesGenerator.new("us-east-1")
      spec.generate(false)
      expect(spec.model_files.count).to be > 300
    end
  end
end