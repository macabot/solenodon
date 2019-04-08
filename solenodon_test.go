package solenodon

import (
	"testing"
	"time"
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
	keys        []interface{}
	dataOut     interface{}
	compareData func(expected, actual interface{}) bool
	nilOut      bool
}

type unmarshal func(b []byte, v interface{}) error

func runGetTests(t *testing.T, tests []*getTest, raw string, unmarshal unmarshal) {
	for i, test := range tests {
		container := &Container{}
		if err := unmarshal([]byte(raw), &container.Data); err != nil {
			t.Errorf("%d, could not unmarshal raw", i)
			continue
		}
		out := container.Get(test.keys...)
		if out == nil {
			if !test.nilOut {
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

type getAndReplaceTest struct {
	keys        []interface{}
	replaceWith interface{}
}

func runGetAndReplaceTests(
	t *testing.T,
	tests []*getAndReplaceTest,
	raw string,
	unmarshal unmarshal,
) {
	for i, test := range tests {
		container := &Container{}
		if err := unmarshal([]byte(raw), &container.Data); err != nil {
			t.Errorf("%d, could decode data", i)
			continue
		}
		if out := container.Get(test.keys...).Replace(test.replaceWith); out == nil {
			t.Errorf("%d, unexpected nil container after replace", i)
			continue
		}
		out := container.Get(test.keys...)
		if out == nil {
			t.Errorf("%d, unexpected nil container at confirmation", i)
			continue
		}
		if out.Data != test.replaceWith {
			t.Errorf("%d, expected data after replacement '%v' (%T), got '%v' (%T)",
				i, test.replaceWith, test.replaceWith, out.Data, out.Data)
		}
	}
}

type deleteTest struct {
	keys []interface{}
}

func runDeleteTests(t *testing.T, tests []*deleteTest, raw string, unmarshal unmarshal) {
	for i, test := range tests {
		container := &Container{}
		if err := unmarshal([]byte(raw), &container.Data); err != nil {
			t.Errorf("%d, could decode data", i)
			continue
		}
		if out := container.Get(test.keys...); out == nil {
			t.Errorf("%d, value not found before delete", i)
			continue
		}
		if out := container.Delete(test.keys...).Get(test.keys...); out != nil {
			t.Errorf("%d, value is found after delete", i)
		}
	}
}

func TestGetFromNilContainerReturnsSelf(t *testing.T) {
	var container *Container
	out := container.Get()
	if container != out {
		t.Error("expected get on nil container to return itself")
	}
}

func TestGetFromUnknownDataType(t *testing.T) {
	container := &Container{Data: foo{}}
	if container.Get("a") != nil {
		t.Error("expected nil on getting from unknown data type")
	}
}

func TestReplaceFromNilContainerReturnsSelf(t *testing.T) {
	var container *Container
	if container.Replace("foo") != container {
		t.Error("expected repalce on nil container to return itself")
	}
}

func TestReplaceWhenParentIsSliceAndKeyNotInt(t *testing.T) {
	container := &Container{
		Data:   2,
		parent: []interface{}{2, 3},
		key:    "foo",
	}
	if container.Replace(22) != nil {
		t.Error("expected nil when replacing with parent slice and non-integer key")
	}
}

func TestReplaceWhenParentIsStringMapAndKeyNotSet(t *testing.T) {
	container := &Container{
		Data:   2,
		parent: map[string]interface{}{"foo": 2},
		key:    "bar",
	}
	if container.Replace(22) != nil {
		t.Error("expected nil when replacing with parent string map and key not set")
	}
}

func TestReplaceWhenParentIsStringMapAndKeyNotString(t *testing.T) {
	container := &Container{
		Data:   2,
		parent: map[string]interface{}{"foo": 2},
		key:    2,
	}
	if container.Replace(22) != nil {
		t.Error("expected nil when replacing with parent string map and key not string")
	}
}

func TestReplaceWhenParentIsMapAndKeyNotSet(t *testing.T) {
	container := &Container{
		Data:   2,
		parent: map[interface{}]interface{}{"foo": 2},
		key:    2,
	}
	if container.Replace(22) != nil {
		t.Error("expected nil when replacing with parent map and key not set")
	}
}

func TestReplaceWhenParentHasUnknownType(t *testing.T) {
	container := &Container{
		Data:   2,
		parent: foo{},
		key:    2,
	}
	if container.Replace(22) != nil {
		t.Error("expected nil when replacing with parent of unknown type")
	}
}

func TestDeleteFromNilContainer(t *testing.T) {
	var container *Container
	if container.Delete("foo") != container {
		t.Error("expected delete on nil container to return itself")
	}
}

func TestDeleteWithNoKeys(t *testing.T) {
	container := &Container{Data: 2}
	out := container.Delete()
	if out != container {
		t.Error("expected delete on container with no keys to return itself")
	}
	if out.Data != nil {
		t.Error("expected delete on container with no keys to set data to nil")
	}
}

func TestDeleteParentNotFound(t *testing.T) {
	container := &Container{Data: 2}
	if container.Delete("foo", "bar") != container {
		t.Error("expected delete on container where parent is not found to return itself")
	}
}
