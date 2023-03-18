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

func ResourceWafProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"config": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem:     ResourceWafConfigSchema(),
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"files": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     ResourceWafDataFileSchema(),
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

func resourceAviWafProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviWafProfileCreate,
		Read:   ResourceAviWafProfileRead,
		Update: resourceAviWafProfileUpdate,
		Delete: resourceAviWafProfileDelete,
		Schema: ResourceWafProfileSchema(),
		Importer: &schema.ResourceImporter{
			State: ResourceWafProfileImporter,
		},
	}
}

func ResourceWafProfileImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	s := ResourceWafProfileSchema()
	return ResourceImporter(d, m, "wafprofile", s)
}

func ResourceAviWafProfileRead(d *schema.ResourceData, meta interface{}) error {
	s := ResourceWafProfileSchema()
	err := ApiRead(d, meta, "wafprofile", s)
	if err != nil {
		log.Printf("[ERROR] in reading object %v\n", err)
	}
	return err
}

func resourceAviWafProfileCreate(d *schema.ResourceData, meta interface{}) error {
	s := ResourceWafProfileSchema()
	err := ApiCreateOrUpdate(d, meta, "wafprofile", s)
	if err == nil {
		err = ResourceAviWafProfileRead(d, meta)
	}
	return err
}

func resourceAviWafProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	s := ResourceWafProfileSchema()
	var err error
	err = ApiCreateOrUpdate(d, meta, "wafprofile", s)
	if err == nil {
		err = ResourceAviWafProfileRead(d, meta)
	}
	return err
}

func resourceAviWafProfileDelete(d *schema.ResourceData, meta interface{}) error {
	objType := "wafprofile"
	client := meta.(*clients.AviClient)
	if ApiDeleteSystemDefaultCheck(d) {
		return nil
	}
	uuid := d.Get("uuid").(string)
	if uuid != "" {
		path := "api/" + objType + "/" + uuid
		err := client.AviSession.Delete(path)
		if err != nil && !(strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "204") || strings.Contains(err.Error(), "403")) {
			log.Println("[INFO] resourceAviWafProfileDelete not found")
			return err
		}
		d.SetId("")
	}
	return nil
}
