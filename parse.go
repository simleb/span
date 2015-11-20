package main

import (
	"flag"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

// A Key is a hierachical path to a variable.
type Key []string

// MakeKey parses a string into a Key.
// TODO: handle quoted subkeys
func MakeKey(s string) Key {
	return Key(strings.Split(s, "."))
}

// String returns the string representation of a Key.
// TODO: handle quoted subkeys
func (k Key) String() string {
	return strings.Join([]string(k), ".")
}

// Equals tests if two keys are equal.
func (k Key) Equals(a Key) bool {
	if len(k) != len(a) {
		return false
	}
	for i, v := range k {
		if a[i] != v {
			return false
		}
	}
	return true
}

// KeySet represents a set of Keys.
type KeySet []Key

// String returns the string representation of a KeySet.
func (ks KeySet) String() string {
	var keys []string
	for _, k := range ks {
		keys = append(keys, k.String())
	}
	return strings.Join(keys, ",")
}

// Set parses a string into a KeySet.
// TODO: handle quoted subkeys
func (ks *KeySet) Set(s string) error {
	k := MakeKey(s)
	for _, v := range *ks {
		if v.Equals(k) {
			return nil
		}
	}
	*ks = append(*ks, k)
	return nil
}

// KeySetSet represents a set of KeySets.
type KeySetSet []KeySet

// String returns the string representation of a KeySetSet.
func (kss KeySetSet) String() string {
	var s []string
	for _, ks := range kss {
		s = append(s, ks.String())
	}
	return strings.Join(s, " ")
}

// Set parses a string into a KeySetSet.
func (kss *KeySetSet) Set(s string) error {
	keys := strings.Split(s, ",")
	var ks KeySet
	for _, s := range keys {
		k := MakeKey(s)
		for _, ks := range *kss {
			for _, v := range ks {
				if v.Equals(k) {
					return fmt.Errorf("variable %q is already bound", k)
				}
			}
		}
		ks.Set(s)
	}
	*kss = append(*kss, ks)
	return nil
}

// Find returns the KeySet that contains k or a singleton.
func (kss KeySetSet) Find(k Key) KeySet {
	for _, ks := range kss {
		for _, x := range ks {
			if k.Equals(x) {
				return ks
			}
		}
	}
	return KeySet{k}
}

var (
	renderFlag KeySet
	bindFlag   KeySetSet
)

func init() {
	flag.Var(&renderFlag, "render", "render templates within variable")
	flag.Var(&renderFlag, "r", "short for --render")
	flag.Var(&bindFlag, "bind", "bind variables to be expanded synchronously")
	flag.Var(&bindFlag, "b", "short for --bind")
}

// ParseFlags parses command line flags and returns
// a template for the output path, the input file path,
// the set of keys to be expanded, the set of keys
// to be rendered, the set of bounded keys and an error.
func ParseFlags() (output *template.Template, input string, expand, render KeySet, bind KeySetSet, err error) {
	flag.Parse()

	// check required arguments
	if flag.NArg() < 1 {
		Fatal(fmt.Errorf("output path missing"))
	} else if flag.NArg() < 2 {
		Fatal(fmt.Errorf("input path missing"))
	}
	input = flag.Arg(1)

	// extract name of variables to expand and produce template
	output, expand, err = ParseOutputPath(flag.Arg(0))
	if err != nil {
		return
	}

	// return variables to render and to bind
	render = renderFlag
	bind = bindFlag

	// remove bound variables from expand
	expand = RemoveBound(expand, bind)

	return
}

// RemoveBound filters expand such that bound variables
// don't get expanded separately.
func RemoveBound(expand KeySet, bind KeySetSet) KeySet {
	for i, kx := range expand {
		for _, k := range bind.Find(kx) {
			if k.Equals(kx) {
				continue
			}
			n := len(expand)
			for j := i + 1; j < n; j++ {
				if k.Equals(expand[j]) {
					expand[j] = expand[n-1]
					expand = expand[:n-1]
					break
				}
			}
		}
	}
	return expand
}

// These patterns define the syntax for variable replacement.
var varPattern = regexp.MustCompile(`\{([^|}]+)\}`)
var fmtPattern = regexp.MustCompile(`\{([^|}]+)\|([^}]+)\}`)

// ParseOutputPath extracts variable names from and parses a template output path.
func ParseOutputPath(path string) (*template.Template, KeySet, error) {
	var expand KeySet
	for _, v := range varPattern.FindAllStringSubmatch(path, -1) {
		expand = append(expand, MakeKey(v[1]))
	}
	for _, v := range fmtPattern.FindAllStringSubmatch(path, -1) {
		expand = append(expand, MakeKey(v[1]))
	}
	out, err := ParseTemplate(path)
	if err != nil {
		return nil, nil, fmt.Errorf("output path incorrectly formatted")
	}
	return out, expand, nil
}

// ParseTemplate parses a custom template.
func ParseTemplate(s string) (*template.Template, error) {
	s = varPattern.ReplaceAllString(s, "{{get . `$1`}}")
	s = fmtPattern.ReplaceAllString(s, "{{get . `$1` | printf `$2`}}")
	tpl := template.New("Render variable")
	tpl = tpl.Funcs(template.FuncMap{"get": func(c Config, s string) interface{} {
		return c.Get(MakeKey(s))
	}})
	out, err := tpl.Parse(s)
	if err != nil {
		return nil, err
	}
	return out, nil
}
