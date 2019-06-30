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

module SharedResources

  def resource_tag_property_type
    props = {
      "Documentation" => "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-resource-tags.html",
      "Properties" => {
        "Key" => {
          "Documentation" => "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-resource-tags.html#cfn-resource-tags-key",
          "PrimitiveType" => "String",
          "UpdateType" => "Mutable",
          "Required" => false
        },
        "Value" => {
          "Documentation" => "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-resource-tags.html#cfn-resource-tags-value",
          "PrimitiveType" => "String",
          "UpdateType" => "Mutable",
          "Required" => false
        }
      }
    }
    props
  end

end