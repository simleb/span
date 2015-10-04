// Command span generates the cartesian product of config files.
//
// Command span can be used when you want to generate multiple config files
// with some variables taking different values in each generated file.
// The input is a config file where the value of those variables have been
// replaced with arrays of values, and a path for the output files as a
// template containing reference to these variables. See the example below.
//
// Usage
//
//	span <output file template> <input file>
//
// Example
//
//	span "{mode}/{id|%03d}.toml" config.toml
// where config.toml contains
//	foo = 42
//	id = [5, 10, 200]
//	mode = ["normal", "crazy"]
// will produce the following files:
//	- normal/005.toml
//	- normal/010.toml
//	- normal/200.toml
//	- crazy/005.toml
//	- crazy/010.toml
//	- crazy/200.toml
// And (for instance) crazy/010.toml will contain
//	foo = 42
//	id = 10
//	mode = "crazy"
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"text/template"
)

const usage = `Usage: span <output files template> <input file>

Example: span "{mode}/{id|%03d}.toml" config.toml
`

// A Config represents a config file.
type Config map[string]interface{}

// A Codec encodes and decodes config files.
type Codec interface {
	Encode(w io.Writer, conf Config) error
	Decode(path string, conf *Config) error
}

// codecs is a map of file extensions to Codecs.
var codecs = map[string]Codec{}

func main() {
	if len(os.Args) != 3 {
		Fatal(fmt.Errorf("expected 2 arguments, got %d\n\n%s", len(os.Args)-1, usage))
	}

	// load input config file
	var in Config
	codec, err := FindCodec(os.Args[2])
	if err != nil {
		Fatal(err)
	}
	if err := codec.Decode(os.Args[2], &in); err != nil {
		Fatal(err)
	}

	// parse output files path
	out, vars, err := ParseOutputPath(os.Args[1])
	if err != nil {
		Fatal(err)
	}

	// recursively generate all output files
	if err := GenerateOutput(out, vars, in); err != nil {
		Fatal(err)
	}
}

// Fatal prints an error on the standard output and exits with a non-zero status.
func Fatal(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	os.Exit(1)
}

// These patterns define the syntax for variable replacement.
var varPattern = regexp.MustCompile(`\{([^|}]+)\}`)
var fmtPattern = regexp.MustCompile(`\{([^|}]+)\|([^}]+)\}`)

// ParseOutputPath extracts variable names from and parses a template output path.
func ParseOutputPath(path string) (*template.Template, Config, error) {
	vars := make(Config)
	for _, v := range varPattern.FindAllStringSubmatch(path, -1) {
		vars[v[1]] = nil
	}
	for _, v := range fmtPattern.FindAllStringSubmatch(path, -1) {
		vars[v[1]] = nil
	}
	path = varPattern.ReplaceAllString(path, "{{index . `$1`}}")
	path = fmtPattern.ReplaceAllString(path, "{{index . `$1` | printf `$2`}}")
	out, err := template.New("out").Parse(path)
	if err != nil {
		return nil, nil, fmt.Errorf("incorrect output path template")
	}
	return out, vars, nil
}

// GenerateOutput recursively produces all the output files.
func GenerateOutput(out *template.Template, vars, in Config) error {
	if len(vars) == 0 {
		var buf bytes.Buffer
		if err := out.Execute(&buf, in); err != nil {
			return err
		}
		return WriteConfig(buf.String(), in)
	}
	var k string
	for v := range vars {
		k = v
		break
	}
	delete(vars, k)
	vals, ok := in[k].([]interface{})
	if !ok {
		return fmt.Errorf("%q must be an array", k)
	}
	for _, v := range vals {
		in[k] = v
		if err := GenerateOutput(out, vars, in); err != nil {
			return err
		}
	}
	in[k] = vals
	vars[k] = nil
	return nil
}

// WriteConfig writes a single config file.
func WriteConfig(path string, in Config) (err error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); err == nil {
			err = cerr
		}
	}()

	codec, err := FindCodec(path)
	if err != nil {
		return err
	}
	return codec.Encode(file, in)
}
