package solenodon

import (
	"reflect"
	"testing"

	"github.com/go-yaml/yaml"
)

var rawYAML = `
title: example
owner:
  name: macabot
  time: 2001-02-20T21:03:55Z
database:
  server: 127.0.0.1
  ports:
  - 8080
  - 8080
  - 8081
  threshold: 30.5
  enabled: true
servers:
  alpha:
    ip: 10.0.0.1
  beta:
    ip: 10.0.0.2
    log: false
clients:
  - ["gamma", "delta"]
  - [1, 2]
hosts: ["alpha", "omega"]
`

func TestGetInYAML(t *testing.T) {
	tests := []*getTest{
		{keys: []interface{}{"foo"}, nilOut: true},
		{keys: []interface{}{foo{}}, nilOut: true},
		{keys: []interface{}{"hosts", -1}, nilOut: true},
		{keys: []interface{}{"hosts", 2}, nilOut: true},
		{keys: []interface{}{"hosts", "alpha"}, nilOut: true},
		{keys: []interface{}{"title"}, dataOut: "example"},
		{keys: []interface{}{"owner", "name"}, dataOut: "macabot"},
		{
			keys:    []interface{}{"owner", "time"},
			dataOut: "2001-02-20T21:03:55Z",
		},
		{
			keys:        []interface{}{"database", "ports"},
			dataOut:     []interface{}{8080, 8080, 8081},
			compareData: reflect.DeepEqual,
		},
		{keys: []interface{}{"database", "threshold"}, dataOut: 30.5},
		{keys: []interface{}{"database", "enabled"}, dataOut: true},
		{keys: []interface{}{"servers", "beta", "log"}, dataOut: false},
		{
			keys:        []interface{}{"clients", 0},
			dataOut:     []interface{}{"gamma", "delta"},
			compareData: reflect.DeepEqual,
		},
		{keys: []interface{}{"clients", 1, 1}, dataOut: 2},
	}
	runGetTests(t, tests, rawYAML, yaml.Unmarshal)
}

func TestGetAndReplaceYAML(t *testing.T) {
	tests := []*getAndReplaceTest{
		{
			keys:        []interface{}{"owner"},
			replaceWith: 4,
		},
		{
			keys:        []interface{}{"servers", "alpha", "ip"},
			replaceWith: "123.456.789",
		},
		{
			keys:        []interface{}{"database", "ports", 1},
			replaceWith: 8088,
		},
		{
			keys:        []interface{}{},
			replaceWith: "22",
		},
	}
	runGetAndReplaceTests(t, tests, rawYAML, yaml.Unmarshal)
}

func TestDeleteFromYAML(t *testing.T) {
	tests := []*deleteTest{
		{keys: []interface{}{"servers", "alpha"}},
		{keys: []interface{}{"hosts", 1}},
	}
	runDeleteTests(t, tests, rawYAML, yaml.Unmarshal)
}
