package cli

import (
	"fmt"

	"github.com/bndr/gotabulate"
)

func PrettyPrint(v interface{}) {
	t := gotabulate.Create(v)
	fmt.Println(t.Render("grid"))
}
