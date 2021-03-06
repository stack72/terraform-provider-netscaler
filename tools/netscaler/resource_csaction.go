package netscaler

import (
	"github.com/chiradeep/go-nitro/config/cs"

	"github.com/chiradeep/go-nitro/netscaler"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"fmt"
	"log"
)

func resourceNetScalerCsaction() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Create:        createCsactionFunc,
		Read:          readCsactionFunc,
		Update:        updateCsactionFunc,
		Delete:        deleteCsactionFunc,
		Schema: map[string]*schema.Schema{
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"newname": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"targetlbvserver": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"targetvserverexpr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func createCsactionFunc(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]  netscaler-provider: In createCsactionFunc")
	client := meta.(*NetScalerNitroClient).client
	var csactionName string
	if v, ok := d.GetOk("name"); ok {
		csactionName = v.(string)
	} else {
		csactionName = resource.PrefixedUniqueId("tf-csaction-")
		d.Set("name", csactionName)
	}
	csaction := cs.Csaction{
		Comment:           d.Get("comment").(string),
		Name:              d.Get("name").(string),
		Newname:           d.Get("newname").(string),
		Targetlbvserver:   d.Get("targetlbvserver").(string),
		Targetvserverexpr: d.Get("targetvserverexpr").(string),
	}

	_, err := client.AddResource(netscaler.Csaction.Type(), csactionName, &csaction)
	if err != nil {
		return err
	}

	d.SetId(csactionName)

	err = readCsactionFunc(d, meta)
	if err != nil {
		log.Printf("[ERROR] netscaler-provider: ?? we just created this csaction but we can't read it ?? %s", csactionName)
		return nil
	}
	return nil
}

func readCsactionFunc(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] netscaler-provider:  In readCsactionFunc")
	client := meta.(*NetScalerNitroClient).client
	csactionName := d.Id()
	log.Printf("[DEBUG] netscaler-provider: Reading csaction state %s", csactionName)
	data, err := client.FindResource(netscaler.Csaction.Type(), csactionName)
	if err != nil {
		log.Printf("[WARN] netscaler-provider: Clearing csaction state %s", csactionName)
		d.SetId("")
		return nil
	}
	d.Set("name", data["name"])
	d.Set("comment", data["comment"])
	d.Set("name", data["name"])
	d.Set("newname", data["newname"])
	d.Set("targetlbvserver", data["targetlbvserver"])
	d.Set("targetvserverexpr", data["targetvserverexpr"])

	return nil

}

func updateCsactionFunc(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]  netscaler-provider: In updateCsactionFunc")
	client := meta.(*NetScalerNitroClient).client
	csactionName := d.Get("name").(string)

	csaction := cs.Csaction{
		Name: d.Get("name").(string),
	}
	hasChange := false
	if d.HasChange("comment") {
		log.Printf("[DEBUG]  netscaler-provider: Comment has changed for csaction %s, starting update", csactionName)
		csaction.Comment = d.Get("comment").(string)
		hasChange = true
	}
	if d.HasChange("name") {
		log.Printf("[DEBUG]  netscaler-provider: Name has changed for csaction %s, starting update", csactionName)
		csaction.Name = d.Get("name").(string)
		hasChange = true
	}
	if d.HasChange("newname") {
		log.Printf("[DEBUG]  netscaler-provider: Newname has changed for csaction %s, starting update", csactionName)
		csaction.Newname = d.Get("newname").(string)
		hasChange = true
	}
	if d.HasChange("targetlbvserver") {
		log.Printf("[DEBUG]  netscaler-provider: Targetlbvserver has changed for csaction %s, starting update", csactionName)
		csaction.Targetlbvserver = d.Get("targetlbvserver").(string)
		hasChange = true
	}
	if d.HasChange("targetvserverexpr") {
		log.Printf("[DEBUG]  netscaler-provider: Targetvserverexpr has changed for csaction %s, starting update", csactionName)
		csaction.Targetvserverexpr = d.Get("targetvserverexpr").(string)
		hasChange = true
	}

	if hasChange {
		_, err := client.UpdateResource(netscaler.Csaction.Type(), csactionName, &csaction)
		if err != nil {
			return fmt.Errorf("Error updating csaction %s", csactionName)
		}
	}
	return readCsactionFunc(d, meta)
}

func deleteCsactionFunc(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]  netscaler-provider: In deleteCsactionFunc")
	client := meta.(*NetScalerNitroClient).client
	csactionName := d.Id()
	err := client.DeleteResource(netscaler.Csaction.Type(), csactionName)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
