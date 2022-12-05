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
	headers := []string{"name", "value"}
	rows := make([][]interface{}, 0)
	for key, value := range m {
		rows = append(rows, []interface{}{key, value})
	}
	t := gotabulate.Create(rows)
	t.SetHeaders(headers)
	t.SetAlign("left")
	fmt.Println(t.Render("grid"))
}

func PrintArray(v interface{}) {
	bts, err := json.Marshal(&v)
	if err != nil {
		panic(err)
	}
	marray := make([]map[string]interface{}, 0)
	err = json.Unmarshal(bts, &marray)
	if err != nil {
		panic(err)
	}
	headers := make([]string, 0)
	rows := make([][]interface{}, 0)
	for index, item := range marray {
		row := make([]interface{}, 0)
		for key, value := range item {
			if index == 0 {
				headers = append(headers, key)
			}
			row = append(row, value)
		}
		rows = append(rows, row)
	}
	t := gotabulate.Create(rows)
	t.SetHeaders(headers)
	t.SetAlign("left")
	fmt.Println(t.Render("grid"))
}
