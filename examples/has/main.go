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
