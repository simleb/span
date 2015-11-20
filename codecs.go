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
	Decode(path string, conf *Config) error
}

// codecs is a map of file extensions to Codecs.
var codecs = map[string]Codec{}

func init() {
	// register codecs
	codecs[".toml"] = TOML{}
	codecs[".json"] = JSON{}
}

// FindCodec finds a codec by extension.
func FindCodec(path string) (Codec, error) {
	ext := strings.ToLower(filepath.Ext(path))
	codec, ok := codecs[ext]
	if !ok {
		return nil, fmt.Errorf("unsupported file format %q", ext)
	}
	return codec, nil
}

// TOML is a codec for the TOML format.
type TOML struct{}

// Encode creates a TOML config file.
func (TOML) Encode(w io.Writer, conf Config) error {
	return toml.NewEncoder(w).Encode(conf)
}

// Decode reads a TOML config file.
func (TOML) Decode(path string, conf *Config) error {
	_, err := toml.DecodeFile(path, conf)
	return err
}

// JSON is a codec for the JSON format.
type JSON struct{}

// Encode creates a JSON config file.
func (JSON) Encode(w io.Writer, conf Config) error {
	return json.NewEncoder(w).Encode(conf)
}

// Decode reads a JSON config file.
func (JSON) Decode(path string, conf *Config) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewDecoder(file).Decode(conf)
}
