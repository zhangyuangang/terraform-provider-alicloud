---
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_cluster"
sidebar_current: "docs-alicloud-resource-polardb-cluster"
description: |-
  Provides an PolarDB cluster resource.
---

# alicloud\_polardb\_cluster

Provides an PolarDB cluster resource. A PolarDB cluster is an isolated database
environment in the cloud. A PolarDB cluster can contain multiple user-created
databases.

## Example Usage

### Create a PolarDB MySQL cluster

```
variable "name" {
  default = "polardbClusterconfig"
}
variable "creation" {
  default = "PolarDB"
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
  db_type               = "MySQL"
  db_version            = "5.6"
  db_node_class         = "rds.mysql.s2.large"
  cluster_charge_type   = "Postpaid"
  cluster_name          = "${var.name}"
  vswitch_id            = "${alicloud_vswitch.default.id}"
}
```

## Argument Reference

The following arguments are supported:

* `db_type` - (Required,ForceNew) Database type. Value options: MySQL, Oracle, PostgreSQL.
* `db_version` - (Required,ForceNew) Database version. Value options can refer to the latest docs [CreateDBCluster](https://help.aliyun.com/document_detail/98169.html?spm=a2c4g.11186623.6.1080.34c26267JaBTSL) `DBVersion`.
* `db_node_class` - (Required) The db_node_class of cluster node.
* `cluster_type` - (Required) DB Instance type. For details, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/26312.htm).
* `zone_id` - (ForceNew) The Zone to launch the DB cluster. it supports multiple zone.
* `cluster_network_type` - (Deprecated) If you want to create clusters in VPC network, this parameter must be set.
* `cluster_charge_type` - (Optional) Valid values are `Prepaid`, `Postpaid`, Default to `Postpaid`. Currently, the resource only supports PostPaid to PrePaid.
* `period` - (Optional) The duration that you will buy DB cluster (in month). It is valid when cluster_charge_type is `PrePaid`. Valid values: [1~9], 12, 24, 36. Default to 1.
* `auto_renew_period` - (Optional) Auto-renewal period of an cluster, in the unit of the month. It is valid when cluster_charge_type is `PrePaid`. Valid value:[1~12], Default to 1.
* `renewal_status` - (Optional) Valid values are `AutoRenewal`, `Normal`, `NotRenewal`, Default to `Normal`. 
* `security_ips` - (Optional) List of IP addresses allowed to access all databases of an cluster. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]).
* `vswitch_id` - (ForceNew) The virtual switch ID to launch DB instances in one VPC.
* `maintain_time` - (Deprecated) The maintain_time of cluster.
* `parameters` - (Optional) Set of parameters needs to be set after DB cluster was launched. Available parameters can refer to the latest docs [View database parameter templates](https://www.alibabacloud.com/help/doc-detail/26284.htm) .

-> **NOTE:** Because of data backup and migration, change DB cluster type and storage would cost 15~20 minutes. Please make full preparation before changing them.

## Attributes Reference

The following attributes are exported:

* `id` - The PolarDB cluster ID.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Used when creating the polardb cluster (until it reaches the initial `Running` status). 
* `update` - (Defaults to 30 mins) Used when updating the polardb cluster (until it reaches the initial `Running` status). 
* `delete` - (Defaults to 20 mins) Used when terminating the polardb cluster. 

## Import

PolarDB cluster can be imported using the id, e.g.

```
$ terraform import alicloud_polardb_cluster.example pc-abc12345678
```