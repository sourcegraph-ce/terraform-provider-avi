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

func ResourceTrafficCloneProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"clone_servers": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     ResourceCloneServerSchema(),
		},
		"cloud_ref": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
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
		"preserve_client_ip": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
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

func resourceAviTrafficCloneProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviTrafficCloneProfileCreate,
		Read:   ResourceAviTrafficCloneProfileRead,
		Update: resourceAviTrafficCloneProfileUpdate,
		Delete: resourceAviTrafficCloneProfileDelete,
		Schema: ResourceTrafficCloneProfileSchema(),
		Importer: &schema.ResourceImporter{
			State: ResourceTrafficCloneProfileImporter,
		},
	}
}

func ResourceTrafficCloneProfileImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	s := ResourceTrafficCloneProfileSchema()
	return ResourceImporter(d, m, "trafficcloneprofile", s)
}

func ResourceAviTrafficCloneProfileRead(d *schema.ResourceData, meta interface{}) error {
	s := ResourceTrafficCloneProfileSchema()
	err := ApiRead(d, meta, "trafficcloneprofile", s)
	if err != nil {
		log.Printf("[ERROR] in reading object %v\n", err)
	}
	return err
}

func resourceAviTrafficCloneProfileCreate(d *schema.ResourceData, meta interface{}) error {
	s := ResourceTrafficCloneProfileSchema()
	err := ApiCreateOrUpdate(d, meta, "trafficcloneprofile", s)
	if err == nil {
		err = ResourceAviTrafficCloneProfileRead(d, meta)
	}
	return err
}

func resourceAviTrafficCloneProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	s := ResourceTrafficCloneProfileSchema()
	var err error
	err = ApiCreateOrUpdate(d, meta, "trafficcloneprofile", s)
	if err == nil {
		err = ResourceAviTrafficCloneProfileRead(d, meta)
	}
	return err
}

func resourceAviTrafficCloneProfileDelete(d *schema.ResourceData, meta interface{}) error {
	objType := "trafficcloneprofile"
	client := meta.(*clients.AviClient)
	if ApiDeleteSystemDefaultCheck(d) {
		return nil
	}
	uuid := d.Get("uuid").(string)
	if uuid != "" {
		path := "api/" + objType + "/" + uuid
		err := client.AviSession.Delete(path)
		if err != nil && !(strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "204") || strings.Contains(err.Error(), "403")) {
			log.Println("[INFO] resourceAviTrafficCloneProfileDelete not found")
			return err
		}
		d.SetId("")
	}
	return nil
}
