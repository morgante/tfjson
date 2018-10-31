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
provider "google" {
	region = "us-central1"
  }
  
  resource "google_compute_network" "test" {
	name = "test"
  }
  
  module "inner" {
	source = "./inner"
  }	
`

const innerTF = `
resource "google_compute_subnetwork" "subnet0" {
	name          = "testsub"
	ip_cidr_range = "10.0.0.0/24"
	network       = "test"
  }  
`

const expected = `{
    "destroy": false,
    "google_compute_network.test": {
        "auto_create_subnetworks": "true",
        "destroy": false,
        "destroy_tainted": false,
        "gateway_ipv4": "",
        "id": "",
        "name": "test",
        "project": "",
        "routing_mode": "",
        "self_link": ""
    },
    "inner": {
        "destroy": false,
        "google_compute_subnetwork.subnet0": {
            "creation_timestamp": "",
            "destroy": false,
            "destroy_tainted": false,
            "fingerprint": "",
            "gateway_address": "",
            "id": "",
            "ip_cidr_range": "10.0.0.0/24",
            "name": "testsub",
            "network": "test",
            "project": "",
            "secondary_ip_range.#": "",
            "self_link": ""
        }
    }
}`

const expectedFlat = `{
    "google_compute_network.test": {
        "auto_create_subnetworks": "true",
        "destroy": false,
        "destroy_tainted": false,
        "gateway_ipv4": "",
        "id": "",
        "name": "test",
        "project": "",
        "routing_mode": "",
        "self_link": ""
    },
    "inner.google_compute_subnetwork.subnet0": {
        "creation_timestamp": "",
        "destroy": false,
        "destroy_tainted": false,
        "fingerprint": "",
        "gateway_address": "",
        "id": "",
        "ip_cidr_range": "10.0.0.0/24",
        "name": "testsub",
        "network": "test",
        "project": "",
        "secondary_ip_range.#": "",
        "self_link": ""
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
	run("terraform", "init")
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
