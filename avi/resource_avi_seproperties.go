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

func ResourceSePropertiesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"se_agent_properties": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem:     ResourceSeAgentPropertiesSchema(),
		},
		"se_bootup_properties": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem:     ResourceSeBootupPropertiesSchema(),
		},
		"se_runtime_properties": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem:     ResourceSeRuntimePropertiesSchema(),
		},
		"uuid": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}

func resourceAviSeProperties() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviSePropertiesCreate,
		Read:   ResourceAviSePropertiesRead,
		Update: resourceAviSePropertiesUpdate,
		Delete: resourceAviSePropertiesDelete,
		Schema: ResourceSePropertiesSchema(),
		Importer: &schema.ResourceImporter{
			State: ResourceSePropertiesImporter,
		},
	}
}

func ResourceSePropertiesImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	s := ResourceSePropertiesSchema()
	return ResourceImporter(d, m, "seproperties", s)
}

func ResourceAviSePropertiesRead(d *schema.ResourceData, meta interface{}) error {
	s := ResourceSePropertiesSchema()
	err := ApiRead(d, meta, "seproperties", s)
	if err != nil {
		log.Printf("[ERROR] in reading object %v\n", err)
	}
	return err
}

func resourceAviSePropertiesCreate(d *schema.ResourceData, meta interface{}) error {
	s := ResourceSePropertiesSchema()
	err := ApiCreateOrUpdate(d, meta, "seproperties", s)
	if err == nil {
		err = ResourceAviSePropertiesRead(d, meta)
	}
	return err
}

func resourceAviSePropertiesUpdate(d *schema.ResourceData, meta interface{}) error {
	s := ResourceSePropertiesSchema()
	var err error
	err = ApiCreateOrUpdate(d, meta, "seproperties", s)
	if err == nil {
		err = ResourceAviSePropertiesRead(d, meta)
	}
	return err
}

func resourceAviSePropertiesDelete(d *schema.ResourceData, meta interface{}) error {
	objType := "seproperties"
	client := meta.(*clients.AviClient)
	if ApiDeleteSystemDefaultCheck(d) {
		return nil
	}
	uuid := d.Get("uuid").(string)
	if uuid != "" {
		path := "api/" + objType + "/" + uuid
		err := client.AviSession.Delete(path)
		if err != nil && !(strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "204") || strings.Contains(err.Error(), "403")) {
			log.Println("[INFO] resourceAviSePropertiesDelete not found")
			return err
		}
		d.SetId("")
	}
	return nil
}
