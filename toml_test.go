package solenodon

import (
	"reflect"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
)

var rawTOML = `
title = "example"

clients = [
	["gamma", "delta"],
	[1, 2]
]

hosts = ["alpha", "omega"]

[owner]
name = "macabot"
time = 2001-02-20T21:03:55Z

[database]
server = "127.0.0.1"
ports = [8080, 8080, 8081]
threshold =  30.5
enabled = true

[servers.alpha]
ip = "10.0.0.1"

[servers.beta]
ip = "10.0.0.2"
log = false

[[friends]]
id = 0
name = "Wood Compton"

[[friends]]
id = 1
name = "Nina Andrews"

[[friends]]
id = 2
name = "Catalina Newton"
`

func TestGetInTOML(t *testing.T) {
	tests := []*getTest{
		{keys: []interface{}{"foo"}, nilOut: true},
		{keys: []interface{}{foo{}}, nilOut: true},
		{keys: []interface{}{"hosts", -1}, nilOut: true},
		{keys: []interface{}{"hosts", 2}, nilOut: true},
		{keys: []interface{}{"hosts", "alpha"}, nilOut: true},
		{keys: []interface{}{"friends", 10}, nilOut: true},
		{keys: []interface{}{"title"}, dataOut: "example"},
		{keys: []interface{}{"owner", "name"}, dataOut: "macabot"},
		{
			keys:        []interface{}{"owner", "time"},
			dataOut:     time.Date(2001, 2, 20, 21, 03, 55, 0, time.UTC),
			compareData: equalTimes,
		},
		{
			keys:        []interface{}{"database", "ports"},
			dataOut:     []interface{}{int64(8080), int64(8080), int64(8081)},
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
		{keys: []interface{}{"clients", 1, 1}, dataOut: int64(2)},
	}
	runGetTests(t, tests, rawTOML, toml.Unmarshal)
}

func TestGetAndSetDataTOML(t *testing.T) {
	tests := []*getAndSetDataTest{
		{
			keys:    []interface{}{"owner"},
			setData: 4,
		},
		{
			keys:    []interface{}{"servers", "alpha", "ip"},
			setData: "123.456.789",
		},
		{
			keys:    []interface{}{"database", "ports", 1},
			setData: 8088,
		},
		{
			keys:    []interface{}{},
			setData: "22",
		},
		{
			keys:    []interface{}{"friends", 1},
			setData: "bob",
		},
	}
	runGetAndSetDataTests(t, tests, rawTOML, toml.Unmarshal)
}

func TestDeleteFromTOML(t *testing.T) {
	tests := []*deleteTest{
		{keys: []interface{}{"servers", "alpha"}},
		{keys: []interface{}{"hosts", 1}},
		{keys: []interface{}{"friends", 2}},
	}
	runDeleteTests(t, tests, rawTOML, toml.Unmarshal)
}

func TestHasNameOfSecondFriend(t *testing.T) {
	raw := `
id = "5cab94182baf40b79ef49f1e"
age = 38

[[friends]]
id = 0
name = "Wood Compton"

[[friends]]
id = 1
name = "Nina Andrews"

[[friends]]
id = 2
name = "Catalina Newton"
`
	container, err := NewContainerFromBytes([]byte(raw), toml.Unmarshal)
	if err != nil {
		panic(err)
	}

	if !container.Has("friends", 1, "name") {
		t.Errorf("expected container to have name of second friend")
	}
}
