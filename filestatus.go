package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type Photofile struct {
	FileName string
	FilePath string
	ModTime  time.Time
}

type Photofiles struct {
	Photofiles []*Photofile
}

// main
func main() {

	searchPath := "./img/"
	fis, err := ioutil.ReadDir(searchPath)

	if err != nil {
		panic(err)
	}

	ps := &Photofiles{}

	for _, fi := range fis {


		p  := &Photofile{}

		fullPath := filepath.Join(searchPath, fi.Name())
		//fmt.Println(fullPath)

		f, _ := os.Open(fullPath)
		defer f.Close()

		if file_status, err := f.Stat(); err == nil {
			p.FileName = file_status.Name()
			p.FilePath = fullPath
			p.ModTime = file_status.ModTime()
			ps.Photofiles = append(ps.Photofiles, p)

		}

	}

//構造体にappendした内容を1つずつ出力する
	for i, st := range ps.Photofiles {
		fmt.Print(i)
		fmt.Println(st)
	}
}
