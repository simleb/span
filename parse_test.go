package main

import "testing"

var validKeys = map[string]Key{
	// simple keys
	`a`:     {`a`},
	`a_b`:   {`a_b`},
	`abc`:   {`abc`},
	`a.b`:   {`a`, `b`},
	`a.b.c`: {`a`, `b`, `c`},
	`a.a`:   {`a`, `a`},

	// spaces are allowed and trimed
	` a . b `:    {`a`, `b`},
	"a. \t\r\nb": {`a`, `b`},

	// litterals can be numbers
	`12`:  {`12`},
	`1.2`: {`1`, `2`},

	// complex keys must be quoted
	`"a b"`: {`a b`},
	`"a.b"`: {`a.b`},
	`" a "`: {` a `},
	`"a.\". b.	c"`: {`a.". b.	c`},
	`"a b".c`:       {`a b`, `c`},
	`"a b"."c"`:     {`a b`, `c`},
	`"a b".c."d e"`: {`a b`, `c`, `d e`},

	// unicode
	`охватывать`:   {`охватывать`},
	`"🙂"."😐"."☹️"`: {`🙂`, `😐`, `☹️`},
}

var validKeySets = map[string]KeySet{
	`a`:                      {Key{`a`}},
	`a,b`:                    {Key{`a`}, Key{`b`}},
	`a,b,c`:                  {Key{`a`}, Key{`b`}, Key{`c`}},
	` a , b `:                {Key{`a`}, Key{`b`}},
	`a.b,c `:                 {Key{`a`, `b`}, Key{`c`}},
	`"a b".ab."-ab-",φ."'q"`: {Key{`a b`, `ab`, `-ab-`}, Key{`φ`, `'q`}},
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

var invalidKeySets = []string{
	``,
	`  `,
	`,`,
	`a,`,
	`,a`,
	`a,,b`,
}

func TestParseKey(t *testing.T) {
	for s, k := range validKeys {
		keys, err := ParseKeySet(s)
		if err != nil {
			t.Fatalf("ParseKey(%q): %s", s, err)
		}
		if len(keys) > 1 {
			t.Fatalf("ParseKey(%q): expected single key, got %d keys", s, len(keys))
		}
		key := keys[0]
		if !key.Equals(k) {
			t.Fatalf("ParseKey(%q): expected %q, got %q", s, k, key)
		}
	}
	for _, s := range invalidKeys {
		keys, err := ParseKeySet(s)
		if err == nil {
			t.Fatalf("ParseKey(%q): expected error, got %q", s, keys)
		}
	}
}

func TestParseKeySet(t *testing.T) {
	for s, ks := range validKeySets {
		keys, err := ParseKeySet(s)
		if err != nil {
			t.Fatalf("ParseKey(%q): %s", s, err)
		}
		for i, key := range keys {
			if !key.Equals(ks[i]) {
				t.Fatalf("ParseKey(%q): expected %q, got %q", s, ks, keys)
			}
		}
	}
	for _, s := range invalidKeySets {
		keys, err := ParseKeySet(s)
		if err == nil {
			t.Fatalf("ParseKey(%q): expected error, got %q", s, keys)
		}
	}
}
