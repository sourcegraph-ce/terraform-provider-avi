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

func ResourcePKIProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ca_certs": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     ResourceSSLCertificateSchema(),
		},
		"created_by": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"crl_check": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"crls": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     ResourceCRLSchema(),
		},
		"ignore_peer_chain": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"is_federated": {
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
		"validate_only_leaf_crl": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
	}
}

func resourceAviPKIProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviPKIProfileCreate,
		Read:   ResourceAviPKIProfileRead,
		Update: resourceAviPKIProfileUpdate,
		Delete: resourceAviPKIProfileDelete,
		Schema: ResourcePKIProfileSchema(),
		Importer: &schema.ResourceImporter{
			State: ResourcePKIProfileImporter,
		},
	}
}

func ResourcePKIProfileImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	s := ResourcePKIProfileSchema()
	return ResourceImporter(d, m, "pkiprofile", s)
}

func ResourceAviPKIProfileRead(d *schema.ResourceData, meta interface{}) error {
	s := ResourcePKIProfileSchema()
	err := ApiRead(d, meta, "pkiprofile", s)
	if err != nil {
		log.Printf("[ERROR] in reading object %v\n", err)
	}
	return err
}

func resourceAviPKIProfileCreate(d *schema.ResourceData, meta interface{}) error {
	s := ResourcePKIProfileSchema()
	err := ApiCreateOrUpdate(d, meta, "pkiprofile", s)
	if err == nil {
		err = ResourceAviPKIProfileRead(d, meta)
	}
	return err
}

func resourceAviPKIProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	s := ResourcePKIProfileSchema()
	var err error
	err = ApiCreateOrUpdate(d, meta, "pkiprofile", s)
	if err == nil {
		err = ResourceAviPKIProfileRead(d, meta)
	}
	return err
}

func resourceAviPKIProfileDelete(d *schema.ResourceData, meta interface{}) error {
	objType := "pkiprofile"
	client := meta.(*clients.AviClient)
	if ApiDeleteSystemDefaultCheck(d) {
		return nil
	}
	uuid := d.Get("uuid").(string)
	if uuid != "" {
		path := "api/" + objType + "/" + uuid
		err := client.AviSession.Delete(path)
		if err != nil && !(strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "204") || strings.Contains(err.Error(), "403")) {
			log.Println("[INFO] resourceAviPKIProfileDelete not found")
			return err
		}
		d.SetId("")
	}
	return nil
}
