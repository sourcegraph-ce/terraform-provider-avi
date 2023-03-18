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

func ResourceCloudConnectorUserSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"azure_serviceprincipal": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem:     ResourceAzureServicePrincipalCredentialsSchema(),
		},
		"azure_userpass": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem:     ResourceAzureUserPassCredentialsSchema(),
		},
		"gcp_credentials": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem:     ResourceGCPCredentialsSchema(),
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"nsxt_credentials": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem:     ResourceNsxtCredentialsSchema(),
		},
		"oci_credentials": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem:     ResourceOCICredentialsSchema(),
		},
		"password": {
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			Sensitive:        true,
			DiffSuppressFunc: suppressSensitiveFieldDiffs,
		},
		"private_key": {
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			Sensitive:        true,
			DiffSuppressFunc: suppressSensitiveFieldDiffs,
		},
		"public_key": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"tenant_ref": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"tencent_credentials": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem:     ResourceTencentCredentialsSchema(),
		},
		"uuid": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"vcenter_credentials": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem:     ResourceVCenterCredentialsSchema(),
		},
	}
}

func resourceAviCloudConnectorUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviCloudConnectorUserCreate,
		Read:   ResourceAviCloudConnectorUserRead,
		Update: resourceAviCloudConnectorUserUpdate,
		Delete: resourceAviCloudConnectorUserDelete,
		Schema: ResourceCloudConnectorUserSchema(),
		Importer: &schema.ResourceImporter{
			State: ResourceCloudConnectorUserImporter,
		},
	}
}

func ResourceCloudConnectorUserImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	s := ResourceCloudConnectorUserSchema()
	return ResourceImporter(d, m, "cloudconnectoruser", s)
}

func ResourceAviCloudConnectorUserRead(d *schema.ResourceData, meta interface{}) error {
	s := ResourceCloudConnectorUserSchema()
	err := ApiRead(d, meta, "cloudconnectoruser", s)
	if err != nil {
		log.Printf("[ERROR] in reading object %v\n", err)
	}
	return err
}

func resourceAviCloudConnectorUserCreate(d *schema.ResourceData, meta interface{}) error {
	s := ResourceCloudConnectorUserSchema()
	err := ApiCreateOrUpdate(d, meta, "cloudconnectoruser", s)
	if err == nil {
		err = ResourceAviCloudConnectorUserRead(d, meta)
	}
	return err
}

func resourceAviCloudConnectorUserUpdate(d *schema.ResourceData, meta interface{}) error {
	s := ResourceCloudConnectorUserSchema()
	var err error
	err = ApiCreateOrUpdate(d, meta, "cloudconnectoruser", s)
	if err == nil {
		err = ResourceAviCloudConnectorUserRead(d, meta)
	}
	return err
}

func resourceAviCloudConnectorUserDelete(d *schema.ResourceData, meta interface{}) error {
	objType := "cloudconnectoruser"
	client := meta.(*clients.AviClient)
	if ApiDeleteSystemDefaultCheck(d) {
		return nil
	}
	uuid := d.Get("uuid").(string)
	if uuid != "" {
		path := "api/" + objType + "/" + uuid
		err := client.AviSession.Delete(path)
		if err != nil && !(strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "204") || strings.Contains(err.Error(), "403")) {
			log.Println("[INFO] resourceAviCloudConnectorUserDelete not found")
			return err
		}
		d.SetId("")
	}
	return nil
}
