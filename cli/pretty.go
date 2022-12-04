package cli

import (
	"encoding/json"
	"fmt"

	"github.com/bndr/gotabulate"
)

func PrintObject(v interface{}) {
	bts, err := json.Marshal(&v)
	if err != nil {
		panic(err)
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(bts, &m)
	if err != nil {
		panic(err)
	}
	PrettyPrint(m)
}

func PrettyPrint(v interface{}) {
	t := gotabulate.Create(v)
	fmt.Println(t.Render("grid"))
}
