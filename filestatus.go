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

		fullPath ,_:= filepath.Abs(filepath.Join(searchPath, fi.Name()))
		//fmt.Println(fullPath)

		f, _ := os.Open(fullPath)
		defer f.Close()

		if file_status, err := f.Stat(); err == nil {
			p.FileName = file_status.Name()
			p.FilePath = fullPath

			jst := time.FixedZone("Asia/Tokyo", 9*60*60)

			p.ModTime = file_status.ModTime().In(jst)
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

// //日付の一致確認サンプル
// 	dt := time.Now()
// 	var jst = time.FixedZone("Asia/Tokyo", 9*60*60)
// 	dt1 := time.Date(dt.Year(), dt.Month(), dt.Day(), 0, 0, 0, 0, jst)
// 	dt2 := time.Date(2021, 4, 16, 0, 0, 0, 0, jst)
// 	if dt1.Equal(dt2) {
// 		fmt.Println("match!")
// 	} else {
// 		fmt.Println("unmatch!")
// 	}

}
