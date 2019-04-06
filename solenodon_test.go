package solenodon

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/go-yaml/yaml"
)

type foo struct{}

func equalTimes(expected, actual interface{}) bool {
	e, okE := expected.(time.Time)
	a, okA := actual.(time.Time)
	if !okE || !okA {
		return false
	}
	return e.Equal(a)
}

type getTest struct {
	raw             string
	container       *Container
	keys            []interface{}
	dataOut         interface{}
	compareData     func(expected, actual interface{}) bool
	nilContainerOut bool
}

func runGetTests(t *testing.T, tests []*getTest) {
	for i, test := range tests {
		out := test.container.Get(test.keys...)
		if out == nil {
			if !test.nilContainerOut {
				t.Errorf("%d, unexpected nil container", i)
			}
		} else if test.compareData != nil {
			if !test.compareData(test.dataOut, out.Data) {
				t.Errorf("%d, expected data '%v' (%T), got '%v' (%T)", i, test.dataOut, test.dataOut, out.Data, out.Data)
			}
		} else if out.Data != test.dataOut {
			t.Errorf("%d, expected data '%v' (%T), got '%v' (%T)", i, test.dataOut, test.dataOut, out.Data, out.Data)
		}
	}
}

func TestGetInJSON(t *testing.T) {
	tests := []*getTest{
		{
			raw:     `{"a":{"b":5,"c":[5,4,7,6,3]}}`,
			keys:    []interface{}{"a", "c", 4},
			dataOut: 3.0,
		},
		{
			raw:             `{}`,
			keys:            []interface{}{"foo"},
			nilContainerOut: true,
		},
		{
			raw:             `{"foo": "bar"}`,
			keys:            []interface{}{foo{}},
			nilContainerOut: true,
		},
	}
	for i, test := range tests {
		test.container = &Container{}
		err := json.Unmarshal([]byte(test.raw), &test.container.Data)
		if err != nil {
			t.Errorf("%d, could not json decode raw data: %s", i, err)
		}
	}
	runGetTests(t, tests)
}

func TestGetInYAML(t *testing.T) {
	tests := []*getTest{
		{
			raw:     `{"a":{"b":5,"c":[5,4,7,6,3]}}`,
			keys:    []interface{}{"a", "c", 4},
			dataOut: 3,
		},
		{
			raw:             `{}`,
			keys:            []interface{}{"foo"},
			nilContainerOut: true,
		},
		{
			raw:             `{"foo": "bar"}`,
			keys:            []interface{}{foo{}},
			nilContainerOut: true,
		},
		{
			raw: `
a: Easy!
b:
  c: 2
  d: [3, 4]
`,
			keys:    []interface{}{"a"},
			dataOut: "Easy!",
		},
		{
			raw: `
a: Easy!
b:
  c: 2
  d: [3, 4]
`,
			keys:    []interface{}{"b", "c"},
			dataOut: 2,
		},
		{
			raw: `
a: Easy!
b:
  c: 2
  d: [3, 4]
`,
			keys:        []interface{}{"b", "d"},
			compareData: reflect.DeepEqual,
			dataOut:     []interface{}{3, 4},
		},
		{
			raw: `
a: Easy!
b:
  55: 2
  d: [3, 4]
`,
			keys:    []interface{}{"b", 55},
			dataOut: 2,
		},
	}
	for i, test := range tests {
		test.container = &Container{}
		err := yaml.Unmarshal([]byte(test.raw), &test.container.Data)
		if err != nil {
			t.Errorf("%d, could not yaml decode raw data: %s", i, err)
		}
	}
	runGetTests(t, tests)
}

func TestGetInTOML(t *testing.T) {
	raw := `
# This is a TOML document. Boom.

title = "TOML Example"

[owner]
name = "Tom Preston-Werner"
organization = "GitHub"
bio = "GitHub Cofounder & CEO\nLikes tater tots and beer."
dob = 1979-05-27T07:32:00Z # First class dates? Why not?

[database]
server = "192.168.1.1"
ports = [ 8001, 8001, 8002 ]
connection_max = 5000
enabled = true

[servers]

  # You can indent as you please. Tabs or spaces. TOML don't care.
  [servers.alpha]
  ip = "10.0.0.1"
  dc = "eqdc10"

  [servers.beta]
  ip = "10.0.0.2"
  dc = "eqdc10"

[clients]
data = [ ["gamma", "delta"], [1, 2] ] # just an update to make sure parsers support it

# Line breaks are OK when inside arrays
hosts = [
  "alpha",
  "omega"
]
`
	container := &Container{}
	err := toml.Unmarshal([]byte(raw), &container.Data)
	if err != nil {
		t.Errorf("could not toml decode raw data: %s", err)
		return
	}

	tests := []*getTest{
		{
			container:       container,
			keys:            []interface{}{"foo"},
			nilContainerOut: true,
		},
		{
			container:       container,
			keys:            []interface{}{foo{}},
			nilContainerOut: true,
		},
		{
			container: container,
			keys:      []interface{}{"title"},
			dataOut:   "TOML Example",
		},
		{
			container: container,
			keys:      []interface{}{"owner", "organization"},
			dataOut:   "GitHub",
		},
		{
			container:   container,
			keys:        []interface{}{"owner", "dob"},
			dataOut:     time.Date(1979, 5, 27, 7, 32, 0, 0, time.UTC),
			compareData: equalTimes,
		},
		{
			container:   container,
			keys:        []interface{}{"database", "ports"},
			dataOut:     []interface{}{int64(8001), int64(8001), int64(8002)},
			compareData: reflect.DeepEqual,
		},
		{
			container: container,
			keys:      []interface{}{"database", "connection_max"},
			dataOut:   int64(5000),
		},
	}
	runGetTests(t, tests)
}

type getAndReplaceTest struct {
	raw       string
	unmarshal func(b []byte, v interface{}) error
	keys      []interface{}
	// replace   func(container *Container)
	replaceWith interface{}
	marshal     func(v interface{}) ([]byte, error)
	out         string
}

func runGetAndReplaceTests(t *testing.T, tests []*getAndReplaceTest) {
	for i, test := range tests {
		container := &Container{}
		err := test.unmarshal([]byte(test.raw), &container.Data)
		if err != nil {
			t.Errorf("%d, could decode data", i)
			continue
		}
		out := container.Get(test.keys...).Replace(test.replaceWith)
		if out == nil {
			t.Errorf("%d, unexpected nil container", i)
			continue
		}
		b, err := test.marshal(container.Data)
		if err != nil {
			t.Errorf("%d, could not encode data", i)
			continue
		}
		s := string(b)
		if s != test.out {
			t.Errorf("%d, expected encoded data '%s', got '%s'", i, test.out, s)
		}
	}
}

func TestGetAndReplaceJSON(t *testing.T) {
	raw := `
{
	"a": {
		"b": [5,2,3],
		"c": 5
	},
	"d": 4.4,
	"e": [
		{"x": "hi"},
		"ho"
	]
}`
	tests := []*getAndReplaceTest{
		{
			raw:         raw,
			unmarshal:   json.Unmarshal,
			keys:        []interface{}{"a"},
			replaceWith: 4,
			marshal:     json.Marshal,
			out:         `{"a":4,"d":4.4,"e":[{"x":"hi"},"ho"]}`,
		},
		{
			raw:         raw,
			unmarshal:   json.Unmarshal,
			keys:        []interface{}{"a", "b", 1},
			replaceWith: 4,
			marshal:     json.Marshal,
			out:         `{"a":{"b":[5,4,3],"c":5},"d":4.4,"e":[{"x":"hi"},"ho"]}`,
		},
		{
			raw:         raw,
			unmarshal:   json.Unmarshal,
			keys:        []interface{}{"e", 0, "x"},
			replaceWith: "bar",
			marshal:     json.Marshal,
			out:         `{"a":{"b":[5,2,3],"c":5},"d":4.4,"e":[{"x":"bar"},"ho"]}`,
		},
		{
			raw:         raw,
			unmarshal:   json.Unmarshal,
			keys:        []interface{}{},
			replaceWith: "22",
			marshal:     json.Marshal,
			out:         `"22"`,
		},
	}
	runGetAndReplaceTests(t, tests)
}
