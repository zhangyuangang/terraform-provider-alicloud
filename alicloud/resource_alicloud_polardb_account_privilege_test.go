package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"testing"
)

func testAccCheckPolarDBAccountPrivilegeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_polardb_account_privilege" {
			continue
		}
		request := polardb.CreateDescribeAccountsRequest()
		request.DBClusterId = rs.Primary.ID
		_, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeAccounts(request)
		})
		if err != nil {
			if IsExceptedError(err, InvalidDBClusterIdNotFound) || IsExceptedError(err, InvalidDBClusterNameNotFound) {
				continue
			}
			return WrapError(err)
		}
	}
	return nil
}

func TestAccAlicloudPolarDBAccountPrivilege_update(t *testing.T) {

	var v *polardb.DBAccount
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sdnsrecordbasic%v.abc", defaultRegionToTest, rand)
	resourceId := "alicloud_polardb_account_privilege.default"
	var basicMap = map[string]string{
		"cluster_id":   CHECKSET,
		"account_name": "tftestprivilege",
		"privilege":    "ReadOnly",
		"db_names.#":   "2",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBAccountPrivilege")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBAccountPrivilegeConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolarDBAccountPrivilegeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":   "${alicloud_polardb_cluster.default.id}",
					"account_name": "${alicloud_polardb_account.default.name}",
					"privilege":    "ReadOnly",
					"db_names":     "${alicloud_polardb_database.default.*.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":   "${alicloud_polardb_cluster.default.id}",
					"account_name": "${alicloud_polardb_account.default.name}",
					"privilege":    "ReadOnly",
					"db_names":     []string{"${alicloud_polardb_database.default.0.name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_names.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":   "${alicloud_polardb_cluster.default.id}",
					"account_name": "${alicloud_polardb_account.default.name}",
					"privilege":    "ReadOnly",
					"db_names":     "${alicloud_polardb_database.default.*.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_names.#": "2",
					}),
				),
			},
		},
	})

}

func resourcePolarDBAccountPrivilegeConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "PolarDB"
	}

	variable "name" {
		default = "%s"
	}

	variable "instancechargetype" {
		default = "Postpaid"
	}

	variable "engine" {
		default = "MySQL"
	}

	variable "engineversion" {
		default = "8.0"
	}

	variable "instanceclass" {
		default = "polar.mysql.x4.large"
	}

	resource "alicloud_polardb_cluster" "default" {
		db_type = "${var.engine}"
		db_version = "${var.engineversion}"
		cluster_charge_type = "${var.instancechargetype}"
		db_node_class = "${var.instanceclass}"
		vswitch_id = "${alicloud_vswitch.default.id}"
		cluster_name = "${var.name}"
	}
	resource "alicloud_polardb_database" "default" {
	  count = 2
	  cluster_id = "${alicloud_polardb_cluster.default.id}"
	  name = "tfaccountpri_${count.index}"
	  description = "from terraform"
	}

	resource "alicloud_polardb_account" "default" {
	  cluster_id = "${alicloud_polardb_cluster.default.id}"
	  name = "tftestprivilege"
	  type = "Normal"
	  password = "Test12345"
	  description = "from terraform"
	}
`, PolarDBCommonTestCase, name)
}
