package bigip

import (
	"log"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipNetSelfIP() *schema.Resource {

	return &schema.Resource{
		Create: resourceBigipNetSelfIPCreate,
		Read:   resourceBigipNetSelfIPRead,
		Update: resourceBigipNetSelfIPUpdate,
		Delete: resourceBigipNetSelfIPDelete,
		//Exists: resourceBigipNetSelfIPExists,
		Importer: &schema.ResourceImporter{
			State: resourceBigipNetSelfIPImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the SelfIP",
				//ValidateFunc: validateF5Name,
			},

			"ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SelfIP IP address",
			},

			"vlan": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the vlan",
				//ValidateFunc: validateF5Name,
			},
		},
	}
}

func resourceBigipNetSelfIPCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	ip := d.Get("ip").(string)
	vlan := d.Get("vlan").(string)

	log.Println("[INFO] Creating SelfIP ")

	err := client.CreateSelfIP(name, ip, vlan)
	// err := client.CreateSelfIP(name+"-self", ip, vlan)

	if err != nil {
		return err
	}

	d.SetId(name)

	return resourceBigipNetSelfIPRead(d, meta)
	// return nil
}

func resourceBigipNetSelfIPRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching SelfIP " + name)

	selfIPs, err := client.SelfIPs()
	if err != nil {
		return err
	}
	for _, selfip := range selfIPs.SelfIPs {
		log.Println(selfip.Name)
		if selfip.Name == name {
			return nil
		}
	}

	return nil
}

func resourceBigipNetSelfIPExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching SelfIP " + name)

	selfIPs, err := client.SelfIPs()
	if err != nil {
		return false, err
	}
	for _, selfip := range selfIPs.SelfIPs {
		log.Println(selfip.Name)
		if selfip.Name == name {
			return true, nil
		}
	}

	return false, nil

}

func resourceBigipNetSelfIPUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Updating SelfIP " + name)

	r := &bigip.SelfIP{
		Name:    name,
		Address: d.Get("ip").(string),
		Vlan:    d.Get("vlan").(string),
	}

	return client.ModifySelfIP(name, r)

}

func resourceBigipNetSelfIPDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Println("[INFO] Deleting selfIP " + name)

	return client.DeleteSelfIP(name)
}

func resourceBigipNetSelfIPImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
