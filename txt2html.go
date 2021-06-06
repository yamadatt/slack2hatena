package main

import (
	"fmt"

	"github.com/yosssi/gohtml"
)

func main() {
	h := `
	<b>
						</b>`
	fmt.Println(gohtml.Format(h))
}
