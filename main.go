package main

import (
	"fmt"

	"./tools"
	"github.com/slack-go/slack"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	api := slack.New(os.Getenv("SLACKAPI"))

	searchPath := "./information/"
	fis, err := ioutil.ReadDir(searchPath)

	if err != nil {
		panic(err)
	}

	for _, fi := range fis {

		fullPath := filepath.Join(searchPath, fi.Name())
		fmt.Println(fullPath)

		str := useIoutilReadFile(fullPath)

		rep := regexp.MustCompile(`https?://files.slack.com[\w!\?/\+\-_~=;\.,\*&@#\$%\(\)'\[\]]+`)
		res := rep.FindAllStringSubmatch(str, -1)

		for _, v := range res {
			filename := GetFileNmae(v[0])

			file, err := os.Create(filename)

			// Slackから画像ファイルをダウンロードする
			err = api.GetFile(v[0], file)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}

			// はてなフォトライフに写真をアップロードする
			hatenaimageid := puthatena.WssePut(filename)

			if hatenaimageid != "" {
				str = strings.Replace(str, v[0], hatenaimageid, 1)
			}
			//ダウンロードしたファイルを削除する
			if err := os.Remove(filename); err != nil {
				fmt.Println(err)
			}

			err = ioutil.WriteFile(fullPath, []byte(str), 0664)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
}

// URLを引数にして、URLの最後のファイル名を返す
func GetFileNmae(imgurl string) (filename string) {
	u, err := url.Parse(imgurl)
	if err != nil {
		log.Fatal(err)
	}
	fname := filepath.Base(u.Path)
	return fname
}

// 読み込んだファイルの中身をstringで返す
func useIoutilReadFile(fileName string) (str string) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
