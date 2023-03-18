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

func ResourceAlertSyslogConfigSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"syslog_servers": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     ResourceAlertSyslogServerSchema(),
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

func resourceAviAlertSyslogConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviAlertSyslogConfigCreate,
		Read:   ResourceAviAlertSyslogConfigRead,
		Update: resourceAviAlertSyslogConfigUpdate,
		Delete: resourceAviAlertSyslogConfigDelete,
		Schema: ResourceAlertSyslogConfigSchema(),
		Importer: &schema.ResourceImporter{
			State: ResourceAlertSyslogConfigImporter,
		},
	}
}

func ResourceAlertSyslogConfigImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	s := ResourceAlertSyslogConfigSchema()
	return ResourceImporter(d, m, "alertsyslogconfig", s)
}

func ResourceAviAlertSyslogConfigRead(d *schema.ResourceData, meta interface{}) error {
	s := ResourceAlertSyslogConfigSchema()
	err := ApiRead(d, meta, "alertsyslogconfig", s)
	if err != nil {
		log.Printf("[ERROR] in reading object %v\n", err)
	}
	return err
}

func resourceAviAlertSyslogConfigCreate(d *schema.ResourceData, meta interface{}) error {
	s := ResourceAlertSyslogConfigSchema()
	err := ApiCreateOrUpdate(d, meta, "alertsyslogconfig", s)
	if err == nil {
		err = ResourceAviAlertSyslogConfigRead(d, meta)
	}
	return err
}

func resourceAviAlertSyslogConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	s := ResourceAlertSyslogConfigSchema()
	var err error
	err = ApiCreateOrUpdate(d, meta, "alertsyslogconfig", s)
	if err == nil {
		err = ResourceAviAlertSyslogConfigRead(d, meta)
	}
	return err
}

func resourceAviAlertSyslogConfigDelete(d *schema.ResourceData, meta interface{}) error {
	objType := "alertsyslogconfig"
	client := meta.(*clients.AviClient)
	if ApiDeleteSystemDefaultCheck(d) {
		return nil
	}
	uuid := d.Get("uuid").(string)
	if uuid != "" {
		path := "api/" + objType + "/" + uuid
		err := client.AviSession.Delete(path)
		if err != nil && !(strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "204") || strings.Contains(err.Error(), "403")) {
			log.Println("[INFO] resourceAviAlertSyslogConfigDelete not found")
			return err
		}
		d.SetId("")
	}
	return nil
}
