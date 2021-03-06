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
		container, err := NewContainerFromBytes([]byte(raw), unmarshal)
		if err != nil {
			t.Errorf("%d, could not unmarshal raw", i)
			continue
		}

		out := container.Get(test.keys...)
		if out == nil {
			if !test.nilOut {
				t.Errorf("%d, unexpected nil container", i)
			}
		} else if test.compareData != nil {
			if !test.compareData(test.dataOut, out.Data()) {
				t.Errorf("%d, expected data '%v' (%T), got '%v' (%T)", i, test.dataOut, test.dataOut, out.Data(), out.Data())
			}
		} else if out.Data() != test.dataOut {
			t.Errorf("%d, expected data '%v' (%T), got '%v' (%T)", i, test.dataOut, test.dataOut, out.Data(), out.Data())
		}
	}
}

type getAndSetDataTest struct {
	keys    []interface{}
	setData interface{}
}

func runGetAndSetDataTests(
	t *testing.T,
	tests []*getAndSetDataTest,
	raw string,
	unmarshal unmarshal,
) {
	for i, test := range tests {
		container, err := NewContainerFromBytes([]byte(raw), unmarshal)
		if err != nil {
			t.Errorf("%d, could not unmarshal raw", i)
			continue
		}
		if out := container.Get(test.keys...).SetData(test.setData); out == nil {
			t.Errorf("%d, unexpected nil container after SetData", i)
			continue
		}
		out := container.Get(test.keys...)
		if out == nil {
			t.Errorf("%d, unexpected nil container at confirmation", i)
			continue
		}
		if out.Data() != test.setData {
			t.Errorf("%d, expected data after SetData '%v' (%T), got '%v' (%T)",
				i, test.setData, test.setData, out.Data(), out.Data())
		}
	}
}

type deleteTest struct {
	keys []interface{}
}

func runDeleteTests(t *testing.T, tests []*deleteTest, raw string, unmarshal unmarshal) {
	for i, test := range tests {
		container, err := NewContainerFromBytes([]byte(raw), unmarshal)
		if err != nil {
			t.Errorf("%d, could not unmarshal raw", i)
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
	container := &Container{data: foo{}}
	if container.Get("a") != nil {
		t.Error("expected nil on getting from unknown data type")
	}
}

func TestSetDataFromNilContainerReturnsSelf(t *testing.T) {
	var container *Container
	if container.SetData("foo") != container {
		t.Error("expected SetData on nil container to return itself")
	}
}

func TestSetDataWhenParentIsSliceAndKeyNotInt(t *testing.T) {
	container := &Container{
		data:   2,
		parent: &Container{data: []interface{}{2, 3}},
		key:    "foo",
	}
	if container.SetData(22) != nil {
		t.Error("expected nil when setting data with parent slice and non-integer key")
	}
}

func TestSetDataWhenParentIsStringMapSliceAndKeyNotInt(t *testing.T) {
	container := &Container{
		data:   map[string]interface{}{"foo": "bar"},
		parent: &Container{data: []map[string]interface{}{{"foo": "bar"}}},
		key:    "foo",
	}
	if container.SetData(33) != nil {
		t.Error("expected nil when setting data with parent []map[string]interface{} and non-integer key")
	}
}

func TestSetDataWhenParentIsStringMapAndKeyNotSet(t *testing.T) {
	container := &Container{
		data:   2,
		parent: &Container{data: map[string]interface{}{"foo": 2}},
		key:    "bar",
	}
	if container.SetData(22) != nil {
		t.Error("expected nil when setting data with parent string map and key not set")
	}
}

func TestSetDataWhenParentIsStringMapAndKeyNotString(t *testing.T) {
	container := &Container{
		data:   2,
		parent: &Container{data: map[string]interface{}{"foo": 2}},
		key:    2,
	}
	if container.SetData(22) != nil {
		t.Error("expected nil when setting data with parent string map and key not string")
	}
}

func TestSetDataWhenParentIsMapAndKeyNotSet(t *testing.T) {
	container := &Container{
		data:   2,
		parent: &Container{data: map[interface{}]interface{}{"foo": 2}},
		key:    2,
	}
	if container.SetData(22) != nil {
		t.Error("expected nil when setting data with parent map and key not set")
	}
}

func TestSetDataWhenParentHasUnknownType(t *testing.T) {
	container := &Container{
		data:   2,
		parent: &Container{data: foo{}},
		key:    2,
	}
	if container.SetData(22) != nil {
		t.Error("expected nil when setting data with parent of unknown type")
	}
}

func TestDeleteFromNilContainer(t *testing.T) {
	var container *Container
	if container.Delete("foo") != container {
		t.Error("expected delete on nil container to return itself")
	}
}

func TestDeleteWithNoKeys(t *testing.T) {
	container := &Container{data: 2}
	out := container.Delete()
	if out != container {
		t.Error("expected delete on container with no keys to return itself")
	}
	if out.Data() != nil {
		t.Error("expected delete on container with no keys to set data to nil")
	}
}

func TestDeleteParentNotFound(t *testing.T) {
	container := &Container{data: 2}
	if container.Delete("foo", "bar") != container {
		t.Error("expected delete on container where parent is not found to return itself")
	}
}

func TestGetDataFromNilContainer(t *testing.T) {
	var container *Container
	if container.Data() != nil {
		t.Error("expected nil data from nil container")
	}
}

func TestContainerHas(t *testing.T) {
	container := &Container{data: map[string]interface{}{"foo": "bar"}}
	if !container.Has("foo") {
		t.Error("expected container to have key 'foo'")
	}
}

func TestContainerHasNot(t *testing.T) {
	container := &Container{data: map[string]interface{}{"foo": "bar"}}
	if container.Has("foo", "bar") {
		t.Error("did not expect container to have key 'foo.bar'")
	}
}
