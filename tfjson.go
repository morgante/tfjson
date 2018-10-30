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
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/morgante/tfjson/converter"
)

func main() {
	var flatten bool

	flag.BoolVar(&flatten, "flatten", false, "Specify whether to flatten the JSON output by prepending module names to instance names instead of making the instances child values.")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: tfjson [OPTIONS] planfile.tfplan")
		flag.CommandLine.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	j, err := tfjson(flag.Args()[0], flatten)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(j)
}

func tfjson(planfile string, flatten bool) (string, error) {
	diff, err := converter.ConvertPlan(planfile, flatten)

	if err != nil {
		return "", err
	}

	j, err := json.MarshalIndent(diff, "", "    ")
	if err != nil {
		return "", err
	}

	return string(j), nil
}
