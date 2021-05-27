package main

import (
	"fmt"


	"github.com/slack-go/slack"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"

)

func main() {
	api := slack.New(os.Getenv("SLACKAPI"))	

	getfile := "https://files.slack.com/files-pri/T3APRBZDX-F022S8CH9GS/_______________________________________________________210523-1605__________________________________________________________hevc_nvenc.mp4.jpg"

		
		file, err := os.Create("file-test.jpg")

		err = api.GetFile(getfile, file)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}


}

func GetFileNmae(imgurl string) (filename string) {
// URLを引数にすると、URLの最後のファイル名を返す

	u, err := url.Parse(imgurl)
	if err != nil {
		log.Fatal(err)
	}

	fname := filepath.Base(u.Path)
	return fname

}

func useIoutilReadFile(fileName string) (str string) {
// 読み込んだファイルの中身をstringで返す

	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}
