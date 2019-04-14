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
