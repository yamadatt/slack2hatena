package main

import (
	"fmt"
	"os"

	"github.com/catatsuy/movabletype"
)

func main() {
	entries, _ := movabletype.Parse(os.Stdin)

	for i, e := range entries {
		fmt.Println(i)
		fmt.Println(e.Body)
	}
}