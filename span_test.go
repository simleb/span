package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Example emulates running command span on a sample config file
// with sample flags and checks the correctness of the output.
func Example() {
	tmp, err := ioutil.TempDir("", "span")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmp)

	// create input file
	input := filepath.Join(tmp, "config.toml")
	content := []byte(`id = 42
size = [5, 10, 200]
width = [640, 800]
height = [480, 600]
[simulation]
	mode = ["normal", "crazy"]
[output]
	dir = "data/{width}x{height}_{size|%03d}_{simulation.mode}{id}.dat"`)
	if err := ioutil.WriteFile(input, content, 0777); err != nil {
		panic(err)
	}

	// emulate running the following command:
	//  span -r output.dir -b width,height "config/{simulation.mode}/{width}x{height}_{size|%03d}.toml" config.toml
	args := os.Args
	defer func() { os.Args = args }()
	os.Args = []string{"span", "-r", "output.dir", "-b", "width,height", filepath.Join(tmp, "config", "{simulation.mode}", "{width}x{height}_{size|%03d}.toml"), input}
	main()

	// check output (partially)
	output, err := ioutil.ReadFile(filepath.Join(tmp, "config", "crazy", "640x480_010.toml"))
	if err != nil {
		panic(err)
	}
	expected := []byte(`height = 480
id = 42
size = 10
width = 640

[output]
  dir = "data/640x480_010_crazy42.dat"

[simulation]
  mode = "crazy"
`)
	if bytes.Equal(output, expected) {
		fmt.Println("Success!")
	} else {
		fmt.Println(string(output))
	}

	// Output: Success!
}
