---
layout: "aws"
page_title: "AWS: aws_cognito_user_pool_resource_server"
side_bar_current: "docs-aws-resource-cognito-user-pool-resource-server"
description: |-
  Provides a Cognito User Pool Resource Server resource.
---

# aws_cognito_user_pool_resource_server

Provides a Cognito User Pool Resource Server resource.

## Example Usage

### Basic configuration

```hcl
resource "aws_cognito_user_pool" "example" {
  name 						= "example-pool"
  auto_verified_attributes  = ["email"]
}

resource "aws_cognito_user_pool_resource_server" "example_server" {
  user_pool_id  	= "${aws_cognito_user_pool.example.id}"
  name 	                = "example_name"
  identifier       	= "some-identifier"

  scopes {
    scope_name = "foo"
    scope_description = "bar"
  }

}
```

## Argument Reference

The following arguments are supported:

* `user_pool_id` (Required) - The user pool id
* `name` (Required) - The resource server name
* `identifier` (Required) - The provider identifier. 
* `scopes` (Optional) - A scope consists of scope_name and scope_description.
