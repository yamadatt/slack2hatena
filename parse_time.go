package main

import (
    "fmt"
    "time"
)

func main() {

    layout_blogysnc := "2006-01-02T15:04:05+15:04"
	layout_movabletype := "01/02/2006 03:04:05 PM"

	jst, _ := time.LoadLocation("Asia/Tokyo")

    // パースする対象の文字列です
    // value := "2021-01-15T06:00:00"
	value := "2021-01-15T06:00:00+09:00"
    t, _ := time.ParseInLocation(layout_blogysnc, value,jst)
	fmt.Println(t)
	fmt.Println(t.Format(layout_movabletype))

   


}