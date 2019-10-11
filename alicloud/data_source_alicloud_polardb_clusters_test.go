package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudPolarClustersDataSource(t *testing.T) {
	rand := acctest.RandInt()
	nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_polardb_cluster.default.cluster_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"name_regex": `"^test1234"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_polardb_cluster.default.cluster_name}"`,
			"status":     `"Running"`,
		}),
		fakeConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"name_regex": `"^test1234"`,
			"status":     `"run"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_polardb_cluster.default.cluster_name}"`,
			"db_type":    `"${alicloud_polardb_cluster.default.db_type}"`,
		}),
		fakeConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"name_regex": `"^test1234"`,
			"db_type":    `"Oracle"`,
		}),
	}

	var existPolarClusterMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"names.#":                   "1",
			"clusters.#":                "1",
			"clusters.0.id":             CHECKSET,
			"clusters.0.name":           CHECKSET,
			"clusters.0.charge_type":    "Postpaid",
			"clusters.0.network_type":   "VPC",
			"clusters.0.region_id":      CHECKSET,
			"clusters.0.zone_id":        CHECKSET,
			"clusters.0.expired":        "false",
			"clusters.0.status":         "Running",
			"clusters.0.engine":         "POLARDB",
			"clusters.0.db_type":        "MySQL",
			"clusters.0.db_version":     "8.0",
			"clusters.0.lock_mode":      "Unlock",
			"clusters.0.delete_lock":    "0",
			"clusters.0.create_time":    CHECKSET,
			"clusters.0.vpc_id":         CHECKSET,
			"clusters.0.db_node_number": "2",
			"clusters.0.db_node_class":  "polar.mysql.x4.large",
			"clusters.0.storage_used":   CHECKSET,
			"clusters.0.db_nodes.#":     "2",
			"clusters.0.tags":           CHECKSET,
		}
	}

	var fakePolarClusterMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"clusters.#": "0",
			"ids.#":      "0",
			"names.#":    "0",
		}
	}

	var PolarClusterCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_polardb_clusters.default",
		existMapFunc: existPolarClusterMapFunc,
		fakeMapFunc:  fakePolarClusterMapFunc,
	}

	PolarClusterCheckInfo.dataSourceTestCheck(t, rand, nameConf, idsConf, allConf)
}

func testAccCheckAlicloudPolarClusterDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
	%s
	variable "creation" {
		default = "PolarDB"
	}

	variable "name" {
		default = "pc-testAccDBInstanceConfig_%d"
	}

	resource "alicloud_polardb_cluster" "default" {
		db_type = "MySQL"
		db_version = "8.0"
		cluster_charge_type = "Postpaid"
		db_node_class = "polar.mysql.x4.large"
		vswitch_id = "${alicloud_vswitch.default.id}"
		cluster_name = "${var.name}"
	}
	data "alicloud_polardb_clusters" "default" {
	  %s
	}
`, PolarDBCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}
