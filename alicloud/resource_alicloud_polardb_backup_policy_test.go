package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func testAccCheckPolarDBBackupPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_polardb_backup_policy" {
			continue
		}
		request := polardb.CreateDescribeBackupPolicyRequest()
		request.DBClusterId = rs.Primary.ID
		_, err := client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
			return polardbClient.DescribeBackupPolicy(request)
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

func TestAccAlicloudPolarDBBackupPolicy(t *testing.T) {
	var v *polardb.DescribeBackupPolicyResponse
	resourceId := "alicloud_polardb_backup_policy.default"
	serverFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeBackupPolicy")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccPolarDBbackuppolicy"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBBackupPolicyConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckPolarDBBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":    "${alicloud_polardb_cluster.default.id}",
					"backup_period": []string{"Tuesday", "Wednesday"},
					"backup_time":   "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#":          "2",
						"backup_period.1592931319": "Tuesday",
						"backup_period.1970423419": "Wednesday",
						"backup_time":              "10:00Z-11:00Z",
					}),
				),
			}},
	})
}

func resourcePolarDBBackupPolicyConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "PolarDB"
	}

	variable "name" {
		default = "%s"
	}

	resource "alicloud_polardb_cluster" "default" {
		db_type = "MySQL"
		db_version = "8.0"
		cluster_charge_type = "Postpaid"
		db_node_class = "polar.mysql.x4.large"
		vswitch_id = "${alicloud_vswitch.default.id}"
		cluster_name = "${var.name}"
	}
`, PolarDBCommonTestCase, name)
}
