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

func ResourceStringGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"kv": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     ResourceKeyValueSchema(),
		},
		"labels": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     ResourceKeyValueSchema(),
		},
		"longest_match": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
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
		"type": {
			Type:     schema.TypeString,
			Required: true,
		},
		"uuid": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}

func resourceAviStringGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviStringGroupCreate,
		Read:   ResourceAviStringGroupRead,
		Update: resourceAviStringGroupUpdate,
		Delete: resourceAviStringGroupDelete,
		Schema: ResourceStringGroupSchema(),
		Importer: &schema.ResourceImporter{
			State: ResourceStringGroupImporter,
		},
	}
}

func ResourceStringGroupImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	s := ResourceStringGroupSchema()
	return ResourceImporter(d, m, "stringgroup", s)
}

func ResourceAviStringGroupRead(d *schema.ResourceData, meta interface{}) error {
	s := ResourceStringGroupSchema()
	err := ApiRead(d, meta, "stringgroup", s)
	if err != nil {
		log.Printf("[ERROR] in reading object %v\n", err)
	}
	return err
}

func resourceAviStringGroupCreate(d *schema.ResourceData, meta interface{}) error {
	s := ResourceStringGroupSchema()
	err := ApiCreateOrUpdate(d, meta, "stringgroup", s)
	if err == nil {
		err = ResourceAviStringGroupRead(d, meta)
	}
	return err
}

func resourceAviStringGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	s := ResourceStringGroupSchema()
	var err error
	err = ApiCreateOrUpdate(d, meta, "stringgroup", s)
	if err == nil {
		err = ResourceAviStringGroupRead(d, meta)
	}
	return err
}

func resourceAviStringGroupDelete(d *schema.ResourceData, meta interface{}) error {
	objType := "stringgroup"
	client := meta.(*clients.AviClient)
	if ApiDeleteSystemDefaultCheck(d) {
		return nil
	}
	uuid := d.Get("uuid").(string)
	if uuid != "" {
		path := "api/" + objType + "/" + uuid
		err := client.AviSession.Delete(path)
		if err != nil && !(strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "204") || strings.Contains(err.Error(), "403")) {
			log.Println("[INFO] resourceAviStringGroupDelete not found")
			return err
		}
		d.SetId("")
	}
	return nil
}
