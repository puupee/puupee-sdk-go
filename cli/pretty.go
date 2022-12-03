package cli

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(v interface{}) {
	bts, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(bts))
}
