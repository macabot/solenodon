package main

import (
	"encoding/json"
	"fmt"

	"github.com/macabot/solenodon"
)

func main() {
	raw := `{
  "balance": "€290.77",
  "email": "example@foo.com",
  "tags": [
    "enim",
    "do",
    "occaecat"
  ],
  "friends": [
    {
      "id": 0,
      "name": "Fuller Glass"
    },
    {
      "id": 1,
      "name": "Chambers Quinn"
    },
    {
      "id": 2,
      "name": "Robles Clay"
    }
  ]
}`
	container, err := solenodon.NewContainerFromBytes([]byte(raw), json.Unmarshal)
	if err != nil {
		panic(err)
	}
	container.Get("tags").Replace([]string{"a", "b", "c"})
	container.Get("friends", 0, "name").Replace("Big Bob")
	container.Get("balance").Replace(struct {
		Currency string `json:"currency"`
		Cents    int    `json:"cents"`
	}{Currency: "euro", Cents: 29077})

	b, err := json.MarshalIndent(container.Data(), "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// output:
	// {
	//         "balance": {
	//                 "currency": "euro",
	//                 "cents": 29077
	//         },
	//         "email": "example@foo.com",
	//         "friends": [
	//                 {
	//                         "id": 0,
	//                         "name": "Big Bob"
	//                 },
	//                 {
	//                         "id": 1,
	//                         "name": "Chambers Quinn"
	//                 },
	//                 {
	//                         "id": 2,
	//                         "name": "Robles Clay"
	//                 }
	//         ],
	//         "tags": [
	//                 "a",
	//                 "b",
	//                 "c"
	//         ]
	// }
	//
}
