package aws

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	//"regexp"
	"testing"
)

func TestAccAWSCognitoUserPoolResourceServer_basic(t *testing.T) {
	name := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSCognitoUserPoolResourceServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSCognitoUserPoolResourceServerConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSCognitoUserPoolResourceServerExists("aws_cognito_user_pool_resource_server.basic"),
					//resource.TestMatchResourceAttr("aws_cognito_user_pool_server_resource.name", "arn",
					//	regexp.MustCompile("^arn:aws:cognito-idp:[^:]+:[0-9]{12}:userpool/[\\w-]+_[0-9a-zA-Z]+$")),
					resource.TestCheckResourceAttr("aws_cognito_user_pool_resource_server.basic", "name", "terraform-test-resource-server-"+name),
					resource.TestCheckResourceAttr("aws_cognito_user_pool_resource_server.basic", "identifier", "terraform-test-resource-server-identifier-"+name),
				),
			},
		},
	})
}

func TestAccAWSCognitoUserPoolResourceServer_withScopes(t *testing.T) {
	name := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSCognitoUserPoolResourceServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSCognitoUserPoolResourceServerConfig_withScopes(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSCognitoUserPoolResourceServerExists("aws_cognito_user_pool_resource_server.scopes"),
					//resource.TestMatchResourceAttr("aws_cognito_user_pool_server_resource.name", "arn",
					//	regexp.MustCompile("^arn:aws:cognito-idp:[^:]+:[0-9]{12}:userpool/[\\w-]+_[0-9a-zA-Z]+$")),
					//resource.TestCheckResourceAttr("aws_cognito_user_pool_resource_server.sr", "scopes", "0"),
					resource.TestCheckResourceAttr("aws_cognito_user_pool_resource_server.scopes", "name", "terraform-test-resource-server-"+name),
					resource.TestCheckResourceAttr("aws_cognito_user_pool_resource_server.scopes", "identifier", "terraform-test-resource-server-identifier-"+name),
				),
			},
		},
	})
}

func testAccAWSCognitoUserPoolResourceServerConfig_basic(name string) string {
	return fmt.Sprintf(`

resource "aws_cognito_user_pool" "pool" {
  name = "terraform-test-pool-%s"
}

resource "aws_cognito_user_pool_resource_server" "basic" {
  name = "terraform-test-resource-server-%s"
  identifier = "terraform-test-resource-server-identifier-%s"
  user_pool_id = "${aws_cognito_user_pool.pool.id}"
}`, name, name, name)
}

func testAccAWSCognitoUserPoolResourceServerConfig_withScopes(name string) string {
	return fmt.Sprintf(`

resource "aws_cognito_user_pool" "pool" {
  name = "terraform-test-pool-%s"
}

resource "aws_cognito_user_pool_resource_server" "scopes" {
  name = "terraform-test-resource-server-%s"
  identifier = "terraform-test-resource-server-identifier-%s"
  user_pool_id = "${aws_cognito_user_pool.pool.id}"
  scopes { 
      scope_name = "foo"
      scope_description = "bar" 
  }

}`, name, name, name)
}

func testAccCheckAWSCognitoUserPoolResourceServerDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).cognitoidpconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_cognito_user_pool_resource_server" {
			continue
		}

		params := &cognitoidentityprovider.DescribeResourceServerInput{
			Identifier: aws.String(rs.Primary.ID),
			UserPoolId: aws.String(rs.Primary.Attributes["user_pool_id"]),
		}

		_, err := conn.DescribeResourceServer(params)

		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "ResourceNotFoundException" {
				return nil
			}
			return err
		}
	}

	return nil
}

func testAccCheckAWSCognitoUserPoolResourceServerExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Cognito Resource Server ID set")
		}

		conn := testAccProvider.Meta().(*AWSClient).cognitoidpconn

		params := &cognitoidentityprovider.DescribeResourceServerInput{
			Identifier: aws.String(rs.Primary.ID),
			UserPoolId: aws.String(rs.Primary.Attributes["user_pool_id"]),
		}

		_, err := conn.DescribeResourceServer(params)

		if err != nil {
			return err
		}

		return nil
	}
}
