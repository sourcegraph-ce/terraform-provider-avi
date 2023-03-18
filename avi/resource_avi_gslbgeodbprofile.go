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

func ResourceGslbGeoDbProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"entries": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     ResourceGslbGeoDbEntrySchema(),
		},
		"is_federated": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
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

func resourceAviGslbGeoDbProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviGslbGeoDbProfileCreate,
		Read:   ResourceAviGslbGeoDbProfileRead,
		Update: resourceAviGslbGeoDbProfileUpdate,
		Delete: resourceAviGslbGeoDbProfileDelete,
		Schema: ResourceGslbGeoDbProfileSchema(),
		Importer: &schema.ResourceImporter{
			State: ResourceGslbGeoDbProfileImporter,
		},
	}
}

func ResourceGslbGeoDbProfileImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	s := ResourceGslbGeoDbProfileSchema()
	return ResourceImporter(d, m, "gslbgeodbprofile", s)
}

func ResourceAviGslbGeoDbProfileRead(d *schema.ResourceData, meta interface{}) error {
	s := ResourceGslbGeoDbProfileSchema()
	err := ApiRead(d, meta, "gslbgeodbprofile", s)
	if err != nil {
		log.Printf("[ERROR] in reading object %v\n", err)
	}
	return err
}

func resourceAviGslbGeoDbProfileCreate(d *schema.ResourceData, meta interface{}) error {
	s := ResourceGslbGeoDbProfileSchema()
	err := ApiCreateOrUpdate(d, meta, "gslbgeodbprofile", s)
	if err == nil {
		err = ResourceAviGslbGeoDbProfileRead(d, meta)
	}
	return err
}

func resourceAviGslbGeoDbProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	s := ResourceGslbGeoDbProfileSchema()
	var err error
	err = ApiCreateOrUpdate(d, meta, "gslbgeodbprofile", s)
	if err == nil {
		err = ResourceAviGslbGeoDbProfileRead(d, meta)
	}
	return err
}

func resourceAviGslbGeoDbProfileDelete(d *schema.ResourceData, meta interface{}) error {
	objType := "gslbgeodbprofile"
	client := meta.(*clients.AviClient)
	if ApiDeleteSystemDefaultCheck(d) {
		return nil
	}
	uuid := d.Get("uuid").(string)
	if uuid != "" {
		path := "api/" + objType + "/" + uuid
		err := client.AviSession.Delete(path)
		if err != nil && !(strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "204") || strings.Contains(err.Error(), "403")) {
			log.Println("[INFO] resourceAviGslbGeoDbProfileDelete not found")
			return err
		}
		d.SetId("")
	}
	return nil
}
