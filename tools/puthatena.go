package puthatena

//はてなフォトライフに写真をアップロードする
//失敗するときがある（ログインエラー等）ので、リトライするように変更したい

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"net/url"
	"bytes"
	"encoding/xml"
	//	"net/http/httputil" // デバッグでダンプする用
	)

type WSSE struct {
	username  string
	password  string
	nonceSize int
}

type Entry struct {
	XMLName xml.Name `xml:"entry"`
	Text    string   `xml:",chardata"`
	Xmlns   string   `xml:"xmlns,attr"`
	Hatena  string   `xml:"hatena,attr"`
	Title   string   `xml:"title"`
	Link    []struct {
		Text  string `xml:",chardata"`
		Rel   string `xml:"rel,attr"`
		Type  string `xml:"type,attr"`
		Href  string `xml:"href,attr"`
		Title string `xml:"title,attr"`
	} `xml:"link"`
	Issued string `xml:"issued"`
	Author struct {
		Text string `xml:",chardata"`
		Name string `xml:"name"`
	} `xml:"author"`
	Generator struct {
		Text    string `xml:",chardata"`
		URL     string `xml:"url,attr"`
		Version string `xml:"version,attr"`
	} `xml:"generator"`
	Subject struct {
		Text string `xml:",chardata"`
		Dc   string `xml:"dc,attr"`
	} `xml:"subject"`
	ID             string `xml:"id"`
	Imageurl       string `xml:"imageurl"`
	Imageurlmedium string `xml:"imageurlmedium"`
	Imageurlsmall  string `xml:"imageurlsmall"`
	Syntax         string `xml:"syntax"`
} 

func New(username, password string, nonceSize int) *WSSE {
	if nonceSize < 1 {
		nonceSize = 20
	}
	
	return &WSSE{
		username:  os.Getenv("HATENAUSERNAME"),
		password:  os.Getenv("HATENAPASSWORD"),
		nonceSize: 30,
	}
}

func createNonce(size int) (string, error) {
	nonce := make([]byte, size)

	_, err := rand.Read(nonce)
	if err != nil {
		return "", err
	}


	return fmt.Sprintf("%x", nonce), nil
}

func (w *WSSE) Create() (string, error) {
	created := time.Now().Format(`2006-01-02T15:04:05Z`)

	n, err := createNonce(w.nonceSize)
	if err != nil {
		return "", err
	}

	h := sha1.New()
	h.Write([]byte(n))
	h.Write([]byte(created))
	h.Write([]byte(w.password))
	digest := base64.URLEncoding.EncodeToString(h.Sum(nil))
	nonce := base64.URLEncoding.EncodeToString([]byte(n))

	format := `UsernameToken Username="%s", PasswordDigest="%s", Nonce="%s", Created="%s"`
	return fmt.Sprintf(format, w.username, digest, nonce, created), nil
}


func WssePut(file_name string) (hatena_img_id string) {
	
    file, err := os.Open(file_name)
	if err != nil {
		log.Println("[-]Couldn't opne file: ", file_name)
	}
	defer file.Close()

	//はてなフォトライフにアップロードするファイルをBase64エンコードする
	binary, _ := ioutil.ReadAll(file)
	base64Data := base64.StdEncoding.EncodeToString(binary)

	body := fmt.Sprintf(
	`<entry xmlns="http://purl.org/atom/ns#">
	  <title>%s</title>
	  <content mode="base64" type="image/jpeg">%s</content>
	  <dc:subject>slack</dc:subject>
	</entry>` ,file_name,base64Data)

	v := url.Values{}
	v.Add("entry",body) 
	
	w := New("username", "password", 14)
	
	req, err := http.NewRequest("POST", "https://f.hatena.ne.jp/atom/post", bytes.NewBuffer([]byte(body)))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	s, err := w.Create()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	req.Header.Set("X-WSSE", s)

	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}


	bodyData, err := ioutil.ReadAll(resp.Body)
	
	defer resp.Body.Close()
	if err != nil {
		log.Fatal("Read Err:", err)
		panic(-1)
	}

	bodyStr := string(bodyData)

	hatena_img_id = GetHatenaFormat(bodyStr)


	//fmt.Println(hatena_img_id)

fmt.Println(resp.Status)

	return hatena_img_id
	
}

func GetHatenaFormat(body string)(hatenaformat string){

result := Entry{} 
    if err := xml.Unmarshal([]byte(body), &result); err != nil {
        fmt.Println("XML Unmarshal error:", err)
        return
    }
 
	return fmt.Sprintf("[%s:plain] \n", result.Syntax)

}