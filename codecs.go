package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// A Codec encodes and decodes config files.
type Codec interface {
	Encode(w io.Writer, conf Config) error
	Decode(r io.Reader, conf *Config) error
}

// codecs is a map of file extensions to Codecs.
var codecs = map[string]Codec{
	".toml": TOML{},
	".json": JSON{},
}

// Encode encodes a Config to a file according to its extension.
func Encode(path string, conf Config) error {
	ext := strings.ToLower(filepath.Ext(path))
	codec, ok := codecs[ext]
	if !ok {
		return fmt.Errorf("unsupported file format %q", ext)
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return codec.Encode(file, conf)
}

// Decode decodes a Config from a file according to its extension.
func Decode(path string, conf *Config) error {
	ext := strings.ToLower(filepath.Ext(path))
	codec, ok := codecs[ext]
	if !ok {
		return fmt.Errorf("unsupported file format %q", ext)
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return codec.Decode(file, conf)
}

// TOML is a codec for the TOML format.
type TOML struct{}

// Encode creates a TOML config file.
func (TOML) Encode(w io.Writer, conf Config) error {
	return toml.NewEncoder(w).Encode(conf)
}

// Decode reads a TOML config file.
func (TOML) Decode(r io.Reader, conf *Config) error {
	_, err := toml.DecodeReader(r, conf)
	return err
}

// JSON is a codec for the JSON format.
type JSON struct{}

// Encode creates a JSON config file.
func (JSON) Encode(w io.Writer, conf Config) error {
	return json.NewEncoder(w).Encode(conf)
}

// Decode reads a JSON config file.
func (JSON) Decode(r io.Reader, conf *Config) error {
	return json.NewDecoder(r).Decode(conf)
}
