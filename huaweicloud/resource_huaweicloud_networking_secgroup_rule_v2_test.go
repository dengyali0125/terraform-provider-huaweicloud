package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/security/rules"
)

func TestAccNetworkingV2SecGroupRule_basic(t *testing.T) {
	var secgroupRule rules.SecGroupRule
	var resourceRuleName string = "huaweicloud_networking_secgroup_rule.secgroup_rule_test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupRule_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupRuleExists(resourceRuleName, &secgroupRule),
					resource.TestCheckResourceAttr(resourceRuleName, "direction", "ingress"),
					resource.TestCheckResourceAttr(resourceRuleName, "port_range_min", "80"),
					resource.TestCheckResourceAttr(resourceRuleName, "port_range_max", "80"),
					resource.TestCheckResourceAttr(resourceRuleName, "ethertype", "IPv4"),
					resource.TestCheckResourceAttr(resourceRuleName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceRuleName, "remote_ip_prefix", "0.0.0.0/0"),
				),
			},
			{
				ResourceName:      resourceRuleName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetworkingV2SecGroupRule_remoteGroup(t *testing.T) {
	var secgroupRule rules.SecGroupRule
	var resourceRuleName string = "huaweicloud_networking_secgroup_rule.secgroup_rule_test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupRule_remoteGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupRuleExists(resourceRuleName, &secgroupRule),
					resource.TestCheckResourceAttr(resourceRuleName, "direction", "ingress"),
					resource.TestCheckResourceAttr(resourceRuleName, "port_range_min", "80"),
					resource.TestCheckResourceAttr(resourceRuleName, "port_range_max", "80"),
					resource.TestCheckResourceAttr(resourceRuleName, "protocol", "tcp"),
					resource.TestCheckResourceAttrSet(resourceRuleName, "remote_group_id"),
				),
			},
			{
				ResourceName:      resourceRuleName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetworkingV2SecGroupRule_lowerCaseCIDR(t *testing.T) {
	var secgroupRule rules.SecGroupRule
	var resourceRuleName string = "huaweicloud_networking_secgroup_rule.secgroup_rule_test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupRule_lowerCaseCIDR(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupRuleExists(resourceRuleName, &secgroupRule),
					resource.TestCheckResourceAttr(resourceRuleName, "remote_ip_prefix", "2001:558:fc00::/39"),
				),
			},
			{
				ResourceName:      resourceRuleName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetworkingV2SecGroupRule_numericProtocol(t *testing.T) {
	var secgroupRule rules.SecGroupRule
	var resourceRuleName string = "huaweicloud_networking_secgroup_rule.secgroup_rule_test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupRule_numericProtocol(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupRuleExists(resourceRuleName, &secgroupRule),
					resource.TestCheckResourceAttr(resourceRuleName, "protocol", "6"),
				),
			},
			{
				ResourceName:      resourceRuleName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNetworkingV2SecGroupRuleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_networking_secgroup_rule" {
			continue
		}

		_, err := rules.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Security group rule still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2SecGroupRuleExists(n string, secGroupRule *rules.SecGroupRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.NetworkingV2Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
		}

		found, err := rules.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Security group rule not found")
		}

		*secGroupRule = *found

		return nil
	}
}

func testAccNetworkingV2SecGroupRule_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "secgroup_test" {
  name        = "%s-secgroup"
  description = "terraform security group rule acceptance test"
}
`, rName)
}

func testAccNetworkingV2SecGroupRule_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup_rule" "secgroup_rule_test" {
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_max    = 80
  port_range_min    = 80
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.secgroup_test.id
}
`, testAccNetworkingV2SecGroupRule_base(rName))
}

func testAccNetworkingV2SecGroupRule_remoteGroup(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup_rule" "secgroup_rule_test" {
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_max    = 80
  port_range_min    = 80
  protocol          = "tcp"
  remote_group_id   = huaweicloud_networking_secgroup.secgroup_test.id
  security_group_id = huaweicloud_networking_secgroup.secgroup_test.id
}
`, testAccNetworkingV2SecGroupRule_base(rName))
}

func testAccNetworkingV2SecGroupRule_lowerCaseCIDR(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup_rule" "secgroup_rule_test" {
  direction         = "ingress"
  ethertype         = "IPv6"
  port_range_max    = 80
  port_range_min    = 80
  protocol          = "tcp"
  remote_ip_prefix  = "2001:558:FC00::/39"
  security_group_id = huaweicloud_networking_secgroup.secgroup_test.id
}
`, testAccNetworkingV2SecGroupRule_base(rName))
}

func testAccNetworkingV2SecGroupRule_numericProtocol(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup_rule" "secgroup_rule_test" {
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_max    = 80
  port_range_min    = 80
  protocol          = "6"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.secgroup_test.id
}
`, testAccNetworkingV2SecGroupRule_base(rName))
}
