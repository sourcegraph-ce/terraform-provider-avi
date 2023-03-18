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

func ResourceSystemLimitsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"controller_limits": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem:     ResourceControllerLimitsSchema(),
		},
		"controller_sizes": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     ResourceControllerSizeSchema(),
		},
		"serviceengine_limits": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem:     ResourceServiceEngineLimitsSchema(),
		},
		"uuid": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}

func resourceAviSystemLimits() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviSystemLimitsCreate,
		Read:   ResourceAviSystemLimitsRead,
		Update: resourceAviSystemLimitsUpdate,
		Delete: resourceAviSystemLimitsDelete,
		Schema: ResourceSystemLimitsSchema(),
		Importer: &schema.ResourceImporter{
			State: ResourceSystemLimitsImporter,
		},
	}
}

func ResourceSystemLimitsImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	s := ResourceSystemLimitsSchema()
	return ResourceImporter(d, m, "systemlimits", s)
}

func ResourceAviSystemLimitsRead(d *schema.ResourceData, meta interface{}) error {
	s := ResourceSystemLimitsSchema()
	err := ApiRead(d, meta, "systemlimits", s)
	if err != nil {
		log.Printf("[ERROR] in reading object %v\n", err)
	}
	return err
}

func resourceAviSystemLimitsCreate(d *schema.ResourceData, meta interface{}) error {
	s := ResourceSystemLimitsSchema()
	err := ApiCreateOrUpdate(d, meta, "systemlimits", s)
	if err == nil {
		err = ResourceAviSystemLimitsRead(d, meta)
	}
	return err
}

func resourceAviSystemLimitsUpdate(d *schema.ResourceData, meta interface{}) error {
	s := ResourceSystemLimitsSchema()
	var err error
	err = ApiCreateOrUpdate(d, meta, "systemlimits", s)
	if err == nil {
		err = ResourceAviSystemLimitsRead(d, meta)
	}
	return err
}

func resourceAviSystemLimitsDelete(d *schema.ResourceData, meta interface{}) error {
	objType := "systemlimits"
	client := meta.(*clients.AviClient)
	if ApiDeleteSystemDefaultCheck(d) {
		return nil
	}
	uuid := d.Get("uuid").(string)
	if uuid != "" {
		path := "api/" + objType + "/" + uuid
		err := client.AviSession.Delete(path)
		if err != nil && !(strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "204") || strings.Contains(err.Error(), "403")) {
			log.Println("[INFO] resourceAviSystemLimitsDelete not found")
			return err
		}
		d.SetId("")
	}
	return nil
}
