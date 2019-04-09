# Solenodon

Solenodon is a Go library for dealing with deserialized data for which the structure is dynamic or not known ahead of time.

Solenodon must be combined with a serialization library that is able to deserialize data into a tree structure of maps and slices if the given target is of type `interface{}`.

Supported:
- [encoding/json]
- [github.com/go-yaml/yaml]
- [github.com/BurntSushi/toml]

Unsupported:
- [encoding/xml] - because there is [no standard way](https://groups.google.com/d/msg/golang-nuts/zEmDOp_yFpU/my8RC0K-DQAJ) to map XML to a key-value structure.

[encoding/json]: https://golang.org/pkg/encoding/json/
[github.com/go-yaml/yaml]: github.com/go-yaml/yaml
[github.com/BurntSushi/toml]: github.com/BurntSushi/toml
[encoding/xml]: https://golang.org/pkg/encoding/xml/

## Install

```go
go get github.com/macabot/solenodon
```

## Examples
### Getting values
```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/macabot/solenodon"
)

func main() {
	raw := `{
  "id": "5cab94182baf40b79ef49f1e",
  "age": 38,
  "friends": [
    {
      "id": 0,
      "name": "Wood Compton"
    },
    {
      "id": 1,
      "name": "Nina Andrews"
    },
    {
      "id": 2,
      "name": "Catalina Newton"
    }
  ]
}`
	container, err := solenodon.NewContainerFromBytes([]byte(raw), json.Unmarshal)
	if err != nil {
		panic(err)
	}

	// type assertion is needed to get the data in the desired type
	id, ok := container.Get("id").Data().(string)
	fmt.Println(id, ok) // 5cab94182baf40b79ef49f1e true

	// encoding/json will by default deserialize all number values as float64
	age, ok := container.Get("age").Data().(float64)
	fmt.Println(age, ok) // 38 true

	// indexes of slices should be of type int
	nameOfSecondFriend, ok := container.Get("friends", 1, "name").Data().(string)
	fmt.Println(nameOfSecondFriend, ok) // Nina Andrews true

	// you might not find what you're looking for
	imaginaryFriend, ok := container.Get("friends", 3).Data().(map[string]interface{})
	fmt.Println(imaginaryFriend, ok) // map[] false
}
```

### Has a value
```go
package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/macabot/solenodon"
)

func main() {
	raw := `
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
	container, err := solenodon.NewContainerFromBytes([]byte(raw), toml.Unmarshal)
	if err != nil {
		panic(err)
	}

	hasFriends := container.Has("friends")
	fmt.Println(hasFriends) // true

	hasNameOfSecondFriend := container.Has("friends", 1, "name")
	fmt.Println(hasNameOfSecondFriend) // true

	hasAgeOfSecondFriend := container.Has("friends", 1, "age")
	fmt.Println(hasAgeOfSecondFriend) // false
}
```

### Deleting a value

### Replacing a value