package main

import (
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"os"
	"io/ioutil"
	"path/filepath"
)

// main
func main() {

	searchPath := "./img/"
	fis, err := ioutil.ReadDir(searchPath)

	if err != nil {
		panic(err)
	}

	for _, fi := range fis {

		fullPath := filepath.Join(searchPath, fi.Name())
		fmt.Println(fullPath)


		file, err := os.Open(fullPath)
		if err != nil {
			panic(err)
		}

		x, err := exif.Decode(file)
		if err != nil {
			panic(err)
		}

		tm, _ := x.DateTime()
		fmt.Println("Taken: ", tm)

	}
}
