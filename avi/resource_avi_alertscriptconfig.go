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

func ResourceAlertScriptConfigSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"action_script": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
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

func resourceAviAlertScriptConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviAlertScriptConfigCreate,
		Read:   ResourceAviAlertScriptConfigRead,
		Update: resourceAviAlertScriptConfigUpdate,
		Delete: resourceAviAlertScriptConfigDelete,
		Schema: ResourceAlertScriptConfigSchema(),
		Importer: &schema.ResourceImporter{
			State: ResourceAlertScriptConfigImporter,
		},
	}
}

func ResourceAlertScriptConfigImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	s := ResourceAlertScriptConfigSchema()
	return ResourceImporter(d, m, "alertscriptconfig", s)
}

func ResourceAviAlertScriptConfigRead(d *schema.ResourceData, meta interface{}) error {
	s := ResourceAlertScriptConfigSchema()
	err := ApiRead(d, meta, "alertscriptconfig", s)
	if err != nil {
		log.Printf("[ERROR] in reading object %v\n", err)
	}
	return err
}

func resourceAviAlertScriptConfigCreate(d *schema.ResourceData, meta interface{}) error {
	s := ResourceAlertScriptConfigSchema()
	err := ApiCreateOrUpdate(d, meta, "alertscriptconfig", s)
	if err == nil {
		err = ResourceAviAlertScriptConfigRead(d, meta)
	}
	return err
}

func resourceAviAlertScriptConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	s := ResourceAlertScriptConfigSchema()
	var err error
	err = ApiCreateOrUpdate(d, meta, "alertscriptconfig", s)
	if err == nil {
		err = ResourceAviAlertScriptConfigRead(d, meta)
	}
	return err
}

func resourceAviAlertScriptConfigDelete(d *schema.ResourceData, meta interface{}) error {
	objType := "alertscriptconfig"
	client := meta.(*clients.AviClient)
	if ApiDeleteSystemDefaultCheck(d) {
		return nil
	}
	uuid := d.Get("uuid").(string)
	if uuid != "" {
		path := "api/" + objType + "/" + uuid
		err := client.AviSession.Delete(path)
		if err != nil && !(strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "204") || strings.Contains(err.Error(), "403")) {
			log.Println("[INFO] resourceAviAlertScriptConfigDelete not found")
			return err
		}
		d.SetId("")
	}
	return nil
}
