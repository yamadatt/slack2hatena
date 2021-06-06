package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"
)

// テンプレートへ引き渡すオブジェクト
// メンバはエクスポート(大文字)しておかないと参照できない
type ExportFormat struct {
	Title             string
	Filename          string
	Primarycategory   string
	Secondarycategory string
	Date              string
	Image             string
	Inyou             string
	Body              string
}

func main() {

	//ファイルを開く
	//searchPath := "./test/"
	searchPath := "./dayly/"
	fis, err := ioutil.ReadDir(searchPath)

	if err != nil {
		panic(err)
	}

	// 構造体にアクセスするため、初期化。hはhtmlの略で使用した。

	// 構造体の初期化 

	h := NewExportFormat()

	for _, fi := range fis {
		fullPath := filepath.Join(searchPath, fi.Name())

		// filename
		h.Filename = "general-" + fi.Name()

		//ファイルを一気に読むと複数行の判定が難しいので、1行ずつ読み込む

		fp, err := os.Open(fullPath)
		if err != nil {
			// エラー処理
		}
		defer fp.Close()

		scanner := bufio.NewScanner(fp)
		line_num := 0

		for scanner.Scan() {

			//fmt.Println(line)
			ss := strings.Split(scanner.Text(), ": ")
			line_num++

			if len(ss) <= 1 {
				if line_num > 6 {
					h.Body += ss[0] + "\n"

				}
			}

			//スライスが1以上あったら（スライスが1以上ということは、「：」で分割できてるということ）
			if len(ss) > 1 {
				key, value := ss[0], ss[1]

				switch key {
				case "Title":
					h.Title = value
					break
				case "Date":
					layout_blogysnc := "2006-01-02T15:04:05+15:04"
					layout_movabletype := "01/02/2006 03:04:05 PM"
					jst, _ := time.LoadLocation("Asia/Tokyo")
				
					// パースする対象の文字列です
					// value := "2021-01-15T06:00:00+09:00"
					t, _ := time.ParseInLocation(layout_blogysnc, value,jst)
					


					h.Date = fmt.Sprintf(t.Format(layout_movabletype))
					break
				}
			}

		}

		if err = scanner.Err(); err != nil {
			// エラー処理
		}

		// カテゴリ
		h.Primarycategory = "slack"

		//文字列置換用の関数を呼ぶ
		h.stringsreg()

		//出力用の関数を呼ぶ
		h.make_export()
//
		h = NewExportFormat()
	}
}

func NewExportFormat() *ExportFormat {
	return &ExportFormat{

	}
}

func (h *ExportFormat) make_export() {

	// テンプレートオブジェクトを生成
	tmpl := template.Must(template.ParseFiles("./test.tmpl"))

	// テンプレートへ渡すデータを作る
	g := &ExportFormat{
		h.Title,
		h.Filename,
		h.Primarycategory,
		h.Secondarycategory,
		h.Date,
		h.Image,
		h.Inyou,
		h.Body,
	}

	// テンプレートからテキストを生成して, os.Stdoutへ出力
	errTemplate := tmpl.Execute(os.Stdout, g) // <- 第二引数にデータを渡す
	//ioutil.WriteFile("go-file",errTemplate , os.ModePerm) // <- 第二引数にデータを渡す
	if errTemplate != nil {
		panic(errTemplate)
	}

}

func (h *ExportFormat) stringsreg() {

	ret := h.Body

	r := strings.NewReplacer("よく読まれている記事", "",
		"<br />", "<br/>",
		`color:rgb`, `color: rgb`)
	ret = r.Replace(ret)

	r2 := strings.NewReplacer("color: rgb(102, 102, 102)", "color: gray",
		`color: rgb(255, 0, 0)`, `color: red`,
		`color: rgb(0, 0, 0)`, `color: black`,
		`color: rgb(5, 109, 178)`, `color: blue`,
		`color: rgb(0, 0, 205)`, `color: blue`,
		"<span><br/></span>", "")

	ret = r2.Replace(ret)

	//コメントを正規表現で削除する
	rep := regexp.MustCompile(`<!--[\s\S]*?-->`)
	ret = rep.ReplaceAllString(ret, "")

	//spanの空タグ
	rep3 := regexp.MustCompile(`<span[^>]+?></span>|\sstyle="\s+?"`)
	ret = rep3.ReplaceAllString(ret, "")

	//HTMLを読みやすくするため、ボックス要素のdiv終わりで改行入れる。後で消す。
	rep2 := regexp.MustCompile(`</div>`)
	ret = rep2.ReplaceAllString(ret, "</div>\r\n")

	//改行だけの行を削除する
	rep11 := regexp.MustCompile(`^\n`)
	ret = rep11.ReplaceAllString(ret, "")

	//newresu1.blog.fc2.comが含まれる行を削除する
	rep12 := regexp.MustCompile(`.+http://newresu1.blog.fc2.com.+\n`)
	ret = rep12.ReplaceAllString(ret, "")

	ret = strings.Replace(ret, `<div id="ad_rs" class="ad_rs_c"></div>`, "", -1)
	ret = strings.Replace(ret, `http://kijosokuho.com/archives/`, "http://192.168.1.254/", -1)
	ret = strings.Replace(ret, `<br/><br/>`, "", -1)
	ret = strings.Replace(ret, `<!-- google_AdSense -->`, "", -1)
	ret = strings.Replace(ret, `<!-- garssはここ！evernote -->`, "", -1)

	ret = strings.Replace(ret, `<b></b>`, ``, -1)
	ret = strings.Replace(ret, `<span><br/><b><br/></b></span>`, "", -1)

	ret = strings.Replace(ret, `<span><br/></span>`, ``, -1)

	ret = strings.Replace(ret, `<span><b></b></span>`, ``, -1)

	h.Body = ret

}

func useIoutilReadFile(fileName string) (str string) {
	// 関数説明：読み込んだファイルの中身をstringで返す
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
