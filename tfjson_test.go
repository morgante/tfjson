/*
Copyright (c) 2016 Palantir Technologies

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

const mainTF = `
provider "aws" {
  region = "us-east-1"
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}

module "inner" {
  source = "./inner"
}
`

const innerTF = `
resource "aws_vpc" "inner" {
  cidr_block = "10.0.0.0/8"
}
`

const expected = `{
    "aws_vpc.main": {
        "arn": "",
        "assign_generated_ipv6_cidr_block": "false",
        "cidr_block": "10.0.0.0/16",
        "default_network_acl_id": "",
        "default_route_table_id": "",
        "default_security_group_id": "",
        "destroy": false,
        "destroy_tainted": false,
        "dhcp_options_id": "",
        "enable_classiclink": "",
        "enable_classiclink_dns_support": "",
        "enable_dns_hostnames": "",
        "enable_dns_support": "true",
        "id": "",
        "instance_tenancy": "default",
        "ipv6_association_id": "",
        "ipv6_cidr_block": "",
        "main_route_table_id": ""
    },
    "destroy": false,
    "inner": {
        "aws_vpc.inner": {
            "arn": "",
            "assign_generated_ipv6_cidr_block": "false",
            "cidr_block": "10.0.0.0/8",
            "default_network_acl_id": "",
            "default_route_table_id": "",
            "default_security_group_id": "",
            "destroy": false,
            "destroy_tainted": false,
            "dhcp_options_id": "",
            "enable_classiclink": "",
            "enable_classiclink_dns_support": "",
            "enable_dns_hostnames": "",
            "enable_dns_support": "true",
            "id": "",
            "instance_tenancy": "default",
            "ipv6_association_id": "",
            "ipv6_cidr_block": "",
            "main_route_table_id": ""
        },
        "destroy": false
    }
}`

const expectedFlat = `{
    "aws_vpc.main": {
        "arn": "",
        "assign_generated_ipv6_cidr_block": "false",
        "cidr_block": "10.0.0.0/16",
        "default_network_acl_id": "",
        "default_route_table_id": "",
        "default_security_group_id": "",
        "destroy": false,
        "destroy_tainted": false,
        "dhcp_options_id": "",
        "enable_classiclink": "",
        "enable_classiclink_dns_support": "",
        "enable_dns_hostnames": "",
        "enable_dns_support": "true",
        "id": "",
        "instance_tenancy": "default",
        "ipv6_association_id": "",
        "ipv6_cidr_block": "",
        "main_route_table_id": ""
    },
    "destroy": false,
    "inner.aws_vpc.inner": {
        "arn": "",
        "assign_generated_ipv6_cidr_block": "false",
        "cidr_block": "10.0.0.0/8",
        "default_network_acl_id": "",
        "default_route_table_id": "",
        "default_security_group_id": "",
        "destroy": false,
        "destroy_tainted": false,
        "dhcp_options_id": "",
        "enable_classiclink": "",
        "enable_classiclink_dns_support": "",
        "enable_dns_hostnames": "",
        "enable_dns_support": "true",
        "id": "",
        "instance_tenancy": "default",
        "ipv6_association_id": "",
        "ipv6_cidr_block": "",
        "main_route_table_id": ""
    }
}`

var planPath string

func TestMain(m *testing.M) {

	run := func(name string, arg ...string) {
		if _, err := exec.Command(name, arg...).Output(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				fmt.Println(exitError)
				os.Exit(1)
			} else {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir)

	mainPath := filepath.Join(dir, "main.tf")
	if err := ioutil.WriteFile(mainPath, []byte(mainTF), 0644); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	innerDir := filepath.Join(dir, "inner")
	if err := os.Mkdir(innerDir, 0755); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	innerPath := filepath.Join(innerDir, "main.tf")
	if err := ioutil.WriteFile(innerPath, []byte(innerTF), 0644); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	planPath = filepath.Join(dir, "terraform.tfplan")
	run("terraform", "get", dir)
	run("terraform", "init", dir)
	run("terraform", "plan", "-out="+planPath, dir)

	os.Exit(m.Run())
}
func TestHeirarchical(t *testing.T) {
	j, err := tfjson(planPath, false)
	if err != nil {
		t.Fatal(err)
	}

	if j != expected {
		t.Errorf("Expected: %s\nActual: %s", expected, j)
	}
}

func TestFlat(t *testing.T) {
	j, err := tfjson(planPath, true)
	if err != nil {
		t.Fatal(err)
	}

	if j != expectedFlat {
		t.Errorf("Expected: %s\nActual: %s", expectedFlat, j)
	}
}
