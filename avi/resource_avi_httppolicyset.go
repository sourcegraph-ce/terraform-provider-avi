/*
 * Copyright (c) 2017. Avi Networks.
 * Author: Gaurav Rastogi (grastogi@avinetworks.com)
 *
 */
package avi

import (
	"github.com/avinetworks/sdk/go/clients"
	"github.com/hashicorp/terraform/helper/schema"
	log "github.com/sourcegraph-ce/logrus"
	"strings"
)

func ResourceHTTPPolicySetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cloud_config_cksum": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"created_by": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"http_request_policy": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem:     ResourceHTTPRequestPolicySchema(),
		},
		"http_response_policy": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem:     ResourceHTTPResponsePolicySchema(),
		},
		"http_security_policy": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem:     ResourceHTTPSecurityPolicySchema(),
		},
		"is_internal_policy": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"labels": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     ResourceKeyValueSchema(),
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"tenant_ref": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"uuid": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}

func resourceAviHTTPPolicySet() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviHTTPPolicySetCreate,
		Read:   ResourceAviHTTPPolicySetRead,
		Update: resourceAviHTTPPolicySetUpdate,
		Delete: resourceAviHTTPPolicySetDelete,
		Schema: ResourceHTTPPolicySetSchema(),
		Importer: &schema.ResourceImporter{
			State: ResourceHTTPPolicySetImporter,
		},
	}
}

func ResourceHTTPPolicySetImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	s := ResourceHTTPPolicySetSchema()
	return ResourceImporter(d, m, "httppolicyset", s)
}

func ResourceAviHTTPPolicySetRead(d *schema.ResourceData, meta interface{}) error {
	s := ResourceHTTPPolicySetSchema()
	err := ApiRead(d, meta, "httppolicyset", s)
	if err != nil {
		log.Printf("[ERROR] in reading object %v\n", err)
	}
	return err
}

func resourceAviHTTPPolicySetCreate(d *schema.ResourceData, meta interface{}) error {
	s := ResourceHTTPPolicySetSchema()
	err := ApiCreateOrUpdate(d, meta, "httppolicyset", s)
	if err == nil {
		err = ResourceAviHTTPPolicySetRead(d, meta)
	}
	return err
}

func resourceAviHTTPPolicySetUpdate(d *schema.ResourceData, meta interface{}) error {
	s := ResourceHTTPPolicySetSchema()
	var err error
	err = ApiCreateOrUpdate(d, meta, "httppolicyset", s)
	if err == nil {
		err = ResourceAviHTTPPolicySetRead(d, meta)
	}
	return err
}

func resourceAviHTTPPolicySetDelete(d *schema.ResourceData, meta interface{}) error {
	objType := "httppolicyset"
	client := meta.(*clients.AviClient)
	if ApiDeleteSystemDefaultCheck(d) {
		return nil
	}
	uuid := d.Get("uuid").(string)
	if uuid != "" {
		path := "api/" + objType + "/" + uuid
		err := client.AviSession.Delete(path)
		if err != nil && !(strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "204") || strings.Contains(err.Error(), "403")) {
			log.Println("[INFO] resourceAviHTTPPolicySetDelete not found")
			return err
		}
		d.SetId("")
	}
	return nil
}
