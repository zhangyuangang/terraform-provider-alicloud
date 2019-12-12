package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudPolarClusterEndPointsDataSource(t *testing.T) {
	rand := acctest.RandInt()

	nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarClusterEndPointsDataSourceConfig(rand, map[string]string{
			"cluster_id": `"${alicloud_polardb_cluster.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudPolarClusterEndPointsDataSourceConfig(rand, map[string]string{
			"cluster_id":     `"${alicloud_polardb_cluster.default.id}"`,
			"db_endpoint_id": `"^test1234"`,
		}),
	}

	var existPolarClusterMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"db_endpoints.0.db_endpoint_id": CHECKSET,
		}
	}

	var fakePolarClusterMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"db_endpoints.#": "0",
		}
	}

	var PolarClusterCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_polardb_endpoints.default",
		existMapFunc: existPolarClusterMapFunc,
		fakeMapFunc:  fakePolarClusterMapFunc,
	}

	PolarClusterCheckInfo.dataSourceTestCheck(t, rand, nameConf)
}

func testAccCheckAlicloudPolarClusterEndPointsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`

	variable "creation" {
	  default = "PolarDB"
	}
	
	variable "name" {
		default = "tf-testAccPolarClusterConfig_%d"
	}
	
	data "alicloud_zones" "default" {
	  available_resource_creation = "${var.creation}"
	}
	
	resource "alicloud_vpc" "default" {
	  name       = "${var.name}"
	  cidr_block = "172.16.0.0/16"
	}
	
	resource "alicloud_vswitch" "default" {
	  vpc_id            = "${alicloud_vpc.default.id}"
	  cidr_block        = "172.16.0.0/24"
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	  name              = "${var.name}"
	}
	
	resource "alicloud_polardb_cluster" "default" {
	  db_type           = "MySQL"
	  db_version        = "8.0"
	  db_node_class     = "polar.mysql.x4.large"
	  vswitch_id        = "${alicloud_vswitch.default.id}"
	  description       = "${var.name}"
	}
	
	data "alicloud_polardb_endpoints" "default" {
	  %s
	}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
