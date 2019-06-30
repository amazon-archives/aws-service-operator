dir = File.expand_path('../..', __dir__)
$:.unshift(dir) unless $:.include?(dir)

require "model_files/shame_list"

RSpec.describe ShameList, "RECURSIVE_SHAME_LIST" do
  context "with no changes" do
    it "has three overrides" do
      expect(ShameList::RECURSIVE_SHAME_LIST.count).to eq 3
    end
  end
end