package main

import "testing"

var validKeys = map[string]Key{
	// simple keys
	`a`:     Key{`a`},
	`a_b`:   Key{`a_b`},
	`abc`:   Key{`abc`},
	`a.b`:   Key{`a`, `b`},
	`a.b.c`: Key{`a`, `b`, `c`},
	`a.a`:   Key{`a`, `a`},

	// spaces are allowed and trimed
	` a . b `:    Key{`a`, `b`},
	"a. \t\r\nb": Key{`a`, `b`},

	// litterals can be numbers
	`12`:  Key{`12`},
	`1.2`: Key{`1`, `2`},

	// complex keys must be quoted
	`"a b"`: Key{`a b`},
	`"a.b"`: Key{`a.b`},
	`" a "`: Key{` a `},
	`"a.\". b.	c"`: Key{`a.". b.	c`},
	`"a b".c`:       Key{`a b`, `c`},
	`"a b"."c"`:     Key{`a b`, `c`},
	`"a b".c."d e"`: Key{`a b`, `c`, `d e`},

	// unicode
	`охватывать`:   Key{`охватывать`},
	`"🙂"."😐"."☹️"`: Key{`🙂`, `😐`, `☹️`},
}

var invalidKeys = []string{
	``,
	`  `,
	`.`,
	`.a`,
	`a.`,
	`.a.b`,
	`a.b.`,
	`a..b`,
	`''`,
	`'a'`,
	`'ab'`,
	`""`,
	`a""`,
	`""a`,

	// TODO: replace text/scanner to be able to deal with incomplete strings
	// `"a`
}

func TestParseKey(t *testing.T) {
	for s, k := range validKeys {
		key, err := ParseKey(s)
		if err != nil {
			t.Fatalf("ParseKey(%q): %s", s, err)
		}
		if !key.Equals(k) {
			t.Fatalf("ParseKey(%q): expected %q, got %q", s, k, key)
		}
	}
	for _, s := range invalidKeys {
		key, err := ParseKey(s)
		if err == nil {
			t.Fatalf("ParseKey(%q): expected error, got %q", s, key)
		}
	}
}
