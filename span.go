package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// Fatal prints an error on the standard error stream
// and exits with a non-zero status.
func Fatal(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	os.Exit(1)
}

func main() {
	output, input, expand, render, bind, err := ParseFlags()
	if err != nil {
		Fatal(fmt.Errorf("bad format in output path"))
	}

	// load input config file
	var conf Config
	if err := Decode(input, &conf); err != nil {
		Fatal(err)
	}

	// recursively generate all output files
	if err := GenerateOutput(output, conf, expand, render, bind); err != nil {
		Fatal(err)
	}
}

// GenerateOutput recursively produces all the output files.
func GenerateOutput(out *template.Template, conf Config, expand, render KeySet, bind KeySetSet) error {
	// write output files when all variables have been expanded
	if len(expand) == 0 {
		return WriteConfig(out, conf, render)
	}

	// select a variable to expand
	kx, expand := expand[0], expand[1:]
	vals, ok := conf.Get(kx).([]interface{})
	if !ok {
		return fmt.Errorf("%q must be an array", kx)
	}
	n := len(vals)

	// find bound variables
	ks := bind.Find(kx)
	values := make([][]interface{}, len(ks))
	for i, k := range ks {
		v, ok := conf.Get(k).([]interface{})
		if !ok || len(v) != n {
			return fmt.Errorf("%q must be an array of the same size as %q", k, kx)
		}
		values[i] = v
	}

	// recurse
	for j := range vals {
		for i, k := range ks {
			conf.Set(k, values[i][j])
		}
		if err := GenerateOutput(out, conf, expand, render, bind); err != nil {
			return err
		}
	}
	for i, k := range ks {
		conf.Set(k, values[i])
	}

	return nil
}

// WriteConfig writes a single config file.
func WriteConfig(out *template.Template, conf Config, render KeySet) error {
	// render variables
	mem := make([]string, len(render))
	for i, k := range render {
		mem[i] = conf.Get(k).(string)
		tpl, err := ParseTemplate(mem[i])
		if err != nil {
			return fmt.Errorf("variable %q incorrectly formatted", mem[i])
		}
		var buf bytes.Buffer
		if err := tpl.Execute(&buf, conf); err != nil {
			return err
		}
		conf.Set(k, buf.String())
	}

	// generate file name
	var buf bytes.Buffer
	if err := out.Execute(&buf, conf); err != nil {
		return err
	}
	path := buf.String()

	// create file and subdirectories
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	// encode output file
	if err := Encode(path, conf); err != nil {
		return err
	}

	// restore template values of rendered variables
	for i, k := range render {
		conf.Set(k, mem[i])
	}

	return nil
}

// A Config represents a config file.
type Config map[string]interface{}

// Get returns the value for the given Key in a Config.
func (c Config) Get(k Key) interface{} {
	n := len(k) - 1
	for i := 0; i < n; i++ {
		c = Config(c[k[i]].(map[string]interface{}))
	}
	return c[k[n]]
}

// Set sets the value for the given Key in a Config.
func (c Config) Set(k Key, v interface{}) {
	n := len(k) - 1
	for i := 0; i < n; i++ {
		c = Config(c[k[i]].(map[string]interface{}))
	}
	c[k[n]] = v
}
