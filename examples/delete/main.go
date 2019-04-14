package main

import (
	"fmt"

	"github.com/go-yaml/yaml"
	"github.com/macabot/solenodon"
)

func main() {
	raw := `
hello: world
database:
  enabled: true
  host: localhost
status:
  points: 32
  blobs: [4, 7, 8, 21]
  logs:
    - date: 2019-01-01
      message: foo
    - date: 2019-01-02
      message: bar
`
	container, err := solenodon.NewContainerFromBytes([]byte(raw), yaml.Unmarshal)
	if err != nil {
		panic(err)
	}

	// delete the status logs and the third status blob
	container.Delete("status", "logs").Delete("status", "blobs", 2)

	b, err := yaml.Marshal(container.Data())
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// output:
	// database:
	//   enabled: true
	//   host: localhost
	// hello: world
	// status:
	//   blobs:
	//   - 4
	//   - 7
	//   - 21
	//   points: 32
}
