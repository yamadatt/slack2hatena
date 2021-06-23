package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
	"sort"
)

// ファイルアップロードに必要な構造体
type Photofile struct {
	FileName string
	FilePath string
	ModTime  time.Time
}

//ファイルアップロードに必要な構造体をスライスで扱うための構造体
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

	//構造体の初期化
	ps := &Photofiles{}

	for _, fi := range fis {

		// 構造体の初期化
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
//デバッグ用にインデックスも使用している。


	for i, st := range ps.Photofiles {
		fmt.Print(i)
		fmt.Println(st)
	}


	//日付でソート
	sort.Slice(ps.Photofiles, func(i, j int) bool {
		return ps.Photofiles[i].ModTime.Before(ps.Photofiles[j].ModTime)
	})


	fmt.Println("----------after sort-----------")

	for i, st := range ps.Photofiles {
		fmt.Print(i)
		fmt.Println(st)
	}


}
