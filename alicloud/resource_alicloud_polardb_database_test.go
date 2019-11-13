package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func testAccCheckPolarDBDatabaseDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_polardb_database" {
			continue
		}
		request := polardb.CreateDescribeDatabasesRequest()
		request.DBClusterId = rs.Primary.ID
		_, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeDatabases(request)
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

func TestAccAlicloudPolarDBDatabase_update(t *testing.T) {
	var database *polardb.Database
	resourceId := "alicloud_polardb_database.default"

	//var dbDatabaseBasicMap = map[string]string{
	//	"cluster_id":    CHECKSET,
	//	"name":          "tftestdatabase",
	//	"character_set": "utf8",
	//	"description":   "",
	//}

	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &database, func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePolarDBDatabase")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBdatabase_basic"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBDatabaseConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolarDBDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id": "${alicloud_polardb_cluster.instance.id}",
					"name":       "tftestdatabase",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "from terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"description": "from terraform"}),
				),
			},
		},
	})
}

func resourcePolarDBDatabaseConfigDependence(name string) string {
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

	resource "alicloud_polardb_cluster" "instance" {
		db_type = "${var.engine}"
		db_version = "${var.engineversion}"
		cluster_charge_type = "${var.instancechargetype}"
		db_node_class = "${var.instanceclass}"
		vswitch_id = "${alicloud_vswitch.default.id}"
		cluster_name = "${var.name}"
	}`, PolarDBCommonTestCase, name)
}
