package main

import (
	"fmt"
	"github.com/slack-go/slack"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type Hatena struct {
	Postdate        string
	Messages
}

type Messages struct {
	Message        string
	User           string
	Posttime       string
	PosttimeDetail string
	Text           string
	PostFiles
}

type PostFiles struct {
	FileTitle    string
	FileType     string
	FileMimeType string
	URLPrivate   string
}

type slackmessages []*Hatena

func main() {
	// 現在時刻取得
	now := time.Now()

	// 取得するメッセージのFromとToを決める
	oldest := strconv.FormatInt(time.Date(now.AddDate(0, 0, -1).Year(), now.AddDate(0, 0, -1).Month(), now.AddDate(0, 0, -1).Day(), 6, 0, 0, 0, time.Local).UnixNano(), 10)
	latest := strconv.FormatInt(time.Date(now.AddDate(0, 0, 0).Year(), now.AddDate(0, 0, 0).Month(), now.AddDate(0, 0, 0).Day(), 5, 59, 59, 999999999, time.Local).UnixNano(), 10)

	postdate := fmt.Sprintf(now.AddDate(0, 0, 0).Format("2006-01-02"))

	// Slack の conversation.history 実行
	api := slack.New(os.Getenv("SLACKAPI"))
	param := slack.GetConversationHistoryParameters{
		ChannelID: "CKMJESN6Q",
		Oldest:    oldest[:10] + "." + oldest[11:16],
		Latest:    latest[:10] + "." + latest[11:16],
	}
	history, err := api.GetConversationHistory(&param)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	var m slackmessages

	for i := len(history.Messages) - 1; i >= 0; i-- {

		m = append(m, doMessage(history.Messages[i].Text, history.Messages[i].Username, history.Messages[i].Timestamp, history.Messages[i].Files))

	}

	m = append(m,doPostdate(postdate))

	//テンプレートに出力する
	tmpl := template.Must(template.ParseFiles("hatena.tpl"))
	SaveToFile(tmpl, m, postdate)

}

// doMessage テンプレートに渡すための情報を編集して、構造体に入れる
func doMessage(mes string, user string, timestamp string, files []slack.File) (r *Hatena) {
	r = new(Hatena)

	// slackに投稿したメッセージを格納
	r.Message = mes

	// slackで使用しているunixタイムを馴染みのある表示形式に変更する
	post_time, _ := strconv.ParseInt(fmt.Sprint(timestamp[:10]), 10, 64)
	dtFromUnix := time.Unix(post_time, 0)
	r.Posttime = dtFromUnix.Format("15:04")
	r.PosttimeDetail = dtFromUnix.Format("2006/1/2 15:04:05")
	r.User = user

	var s []string
	for _, fileid := range files {

		s = append(s, fileid.URLPrivate)

	}

	text := fmt.Sprintf("%s", s)

	text = strings.Trim(text, "[|]")

	r.Text = text

	return r
}

func doPostdate(postdate string) (r *Hatena) {
	r = new(Hatena)

	// はてなブログに投稿する日時を格納する
	r.Postdate = postdate
	
	return r
}

// テンプレート機能を使って、マークダウンファイルを作成する
func SaveToFile(tmpl *template.Template, m slackmessages, filename string) {
	nf, err := os.Create(filename + ".md")
	if err != nil {
		log.Println("error createing file", err)
	}
	defer nf.Close()

	err = tmpl.Execute(nf, m)
	if err != nil {
		log.Fatalln(err)
	}

}
