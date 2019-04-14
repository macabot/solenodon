package main

import (
	"encoding/json"
	"fmt"

	"github.com/macabot/solenodon"
)

func main() {
	raw := []byte(`{"foo":"bar","items":[2,3,{"i":6,"j":7}]}`)
	container, err := solenodon.NewContainerFromBytes(raw, json.Unmarshal)
	if err != nil {
		panic(err)
	}
	fmt.Println(container.Has("foo"))                 // true
	fmt.Println(container.Get("foo").Data().(string)) // bar
	container.Get("items", 2, "j").Replace(44)
	container.Delete("items", 0)
	b, err := json.Marshal(container.Data())
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b)) // {"foo":"bar","items":[3,{"i":6,"j":44}]}
}
