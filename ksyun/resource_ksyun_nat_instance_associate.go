/*
Provides a Association of Nat and Instance resource under VPC resource.

# Example Usage

```hcl


data "ksyun_images" "centos-7_5" {
  platform= "centos-7.5"
}
data "ksyun_availability_zones" "default" {
}

resource "ksyun_security_group" "default" {
  vpc_id = "${ksyun_vpc.foo.id}"
  security_group_name="ksyun-security-group-nat"
}

resource "ksyun_instance" "foo" {
  image_id="${data.ksyun_images.centos-7_5.images.0.image_id}"
  instance_type="N3.2B"

  #max_count=1
  #min_count=1
  subnet_id="${ksyun_subnet.foo.id}"
  instance_password="Xuan663222"
  keep_image_login=false
  charge_type="Daily"
  purchase_time=1
  security_group_id=["${ksyun_security_group.default.id}"]
  instance_name="ksyun-kec-tf-nat"
  sriov_net_support="false"
  project_id=100012
}

resource "ksyun_nat" "foo" {
  nat_name = "ksyun-nat-tf"
  nat_mode = "Subnet"
  nat_type = "public"
  band_width = 1
  charge_type = "DailyPaidByTransfer"
  vpc_id = "${ksyun_vpc.foo.id}"
}
resource "ksyun_vpc" "foo" {
	vpc_name        = "tf-vpc-nat"
	cidr_block = "10.0.5.0/24"
}

resource "ksyun_subnet" "foo" {
  subnet_name      = "tf-acc-nat-subnet1"
  cidr_block = "10.0.5.0/24"
  subnet_type = "Normal"
  vpc_id  = "${ksyun_vpc.foo.id}"
  gateway_ip = "10.0.5.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "${data.ksyun_availability_zones.default.availability_zones.0.availability_zone_name}"
}

resource "ksyun_nat_instance_associate" "foo" {
  nat_id = "${ksyun_nat.foo.id}"
  network_interface_id = "${ksyun_instance.foo.network_interface_id}"
}

```

# Import

nat associate can be imported using the `id`, e.g.

```
$ terraform import ksyun_nat_instance_associate.foo $nat_id:$network_interface_id
```
*/
package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunNatInstanceAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunNatInstanceAssociationCreate,
		Read:   resourceKsyunNatInstanceAssociationRead,
		Delete: resourceKsyunNatInstanceAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: importNatAssociate,
		},

		Schema: map[string]*schema.Schema{
			"nat_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the Nat.",
			},
			"network_interface_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "The id of network interface that belong to instance.",
			},
		},
	}
}
func resourceKsyunNatInstanceAssociationCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateNatInstanceAssociate(d, resourceKsyunNatInstanceAssociation())
	if err != nil {
		return fmt.Errorf("error on creating nat instance associate %q, %s", d.Id(), err)
	}
	return resourceKsyunNatInstanceAssociationRead(d, meta)
}

func resourceKsyunNatInstanceAssociationRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetNatAssociate(d, resourceKsyunNatInstanceAssociation())
	if err != nil {
		return fmt.Errorf("error on reading nat instance associate %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunNatInstanceAssociationDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveNatInstanceAssociate(d)
	if err != nil {
		return fmt.Errorf("error on deleting nat instance associate %q, %s", d.Id(), err)
	}
	return err
}
