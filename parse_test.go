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

var validKeySets = map[string]KeySet{
	`a`:                      KeySet{Key{`a`}},
	`a,b`:                    KeySet{Key{`a`}, Key{`b`}},
	`a,b,c`:                  KeySet{Key{`a`}, Key{`b`}, Key{`c`}},
	` a , b `:                KeySet{Key{`a`}, Key{`b`}},
	`a.b,c `:                 KeySet{Key{`a`, `b`}, Key{`c`}},
	`"a b".ab."-ab-",φ."'q"`: KeySet{Key{`a b`, `ab`, `-ab-`}, Key{`φ`, `'q`}},
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
