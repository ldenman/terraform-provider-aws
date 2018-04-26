package aws

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsCognitoUserPoolResourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsCognitoUserPoolResourceServerCreate,
		Read:   resourceAwsCognitoUserPoolResourceServerRead,
		Update: resourceAwsCognitoUserPoolResourceServerUpdate,
		Delete: resourceAwsCognitoUserPoolResourceServerDelete,

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			// "scopes": {
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeString,
			// 	},
			// },

			"identifier": {
				Type:     schema.TypeString,
				Required: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"user_pool_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAwsCognitoUserPoolResourceServerCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).cognitoidpconn
	log.Print("[DEBUG] Creating Cognito Resource Server")

	name := aws.String(d.Get("name").(string))
	params := &cognitoidentityprovider.CreateResourceServerInput{
		Name:       name,
		Identifier: aws.String(d.Get("identifier").(string)),
		UserPoolId: aws.String(d.Get("user_pool_id").(string)),
	}

	// if v, ok := d.GetOk("scopes"); ok {
	// 	params.Scopes = expandStringList(v.([]interface{}))
	// }

	_, err := conn.CreateResourceServer(params)
	if err != nil {
		return fmt.Errorf("Error creating Cognito Resource Server: %s", err)
	}

	d.SetId(*name)

	return resourceAwsCognitoUserPoolResourceServerRead(d, meta)
}

func resourceAwsCognitoUserPoolResourceServerRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).cognitoidpconn
	log.Printf("[DEBUG] Reading Cognito Resource Server: %s", d.Id())

	ret, err := conn.DescribeResourceServer(&cognitoidentityprovider.DescribeResourceServerInput{
		Identifier: aws.String(d.Get("identifier").(string)),
		UserPoolId: aws.String(d.Get("user_pool_id").(string)),
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "ResourceNotFoundException" {
			d.SetId("")
			return nil
		}
		return err
	}

	ip := ret.ResourceServer
	d.Set("name", ip.Name)
	d.Set("user_pool_id", ip.UserPoolId)
	d.Set("identifier", ip.Identifier)

	// if err := d.Set("scopes", flattenStringList(ip.Scopes)); err != nil {
	// 	return fmt.Errorf("[DEBUG] Error setting scopes error: %#v", err)
	// }

	return nil
}

func resourceAwsCognitoUserPoolResourceServerUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).cognitoidpconn
	log.Print("[DEBUG] Updating Cognito Resource Server")

	params := &cognitoidentityprovider.UpdateResourceServerInput{
		UserPoolId: aws.String(d.Get("UserPoolId").(string)),
		Identifier: aws.String(d.Get("Identifier").(string)),
		Name:       aws.String(d.Id()),
	}

	// if d.HasChange("scopes") {
	// 	params.Scopes = expandStringList(d.Get("scopes").([]interface{}))
	// }

	_, err := conn.UpdateResourceServer(params)
	if err != nil {
		return fmt.Errorf("Error updating Cognito Resource Server: %s", err)
	}

	return resourceAwsCognitoUserPoolResourceServerRead(d, meta)
}

func resourceAwsCognitoUserPoolResourceServerDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).cognitoidpconn
	log.Printf("[DEBUG] Deleting Cognito Resource Server: %s", d.Id())

	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := conn.DeleteResourceServer(&cognitoidentityprovider.DeleteResourceServerInput{
			Identifier: aws.String(d.Get("identifier").(string)),
			UserPoolId: aws.String(d.Get("user_pool_id").(string)),
		})

		if err == nil {
			d.SetId("")
			return nil
		}

		return resource.NonRetryableError(err)
	})
}
