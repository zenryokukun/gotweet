package gotweet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dghubble/oauth1"
)

const (
	// tweet用URL　v2対応
	URL_TWEET = "https://api.twitter.com/2/tweets"
	// 画像アップロード用URL V2未リリース（v1.1を使えとのこと）
	URL_IMAGE = "https://upload.twitter.com/1.1/media/upload.json"
	// twitter認証　ファイル名
	CREDENTIALS = "twitter_conf.json"
)

type (
	keys struct {
		API_KEY       string `json:"API_KEY"`
		API_SECRET    string `json:"API_SECRET"`
		ACCESS_TOKEN  string `json:"ACCESS_TOKEN"`
		ACCESS_SECRET string `json:"ACCESS_SECRET"`
	}

	// uploadのレスポンスの一部
	image struct {
		Image_type string `json:"image_type"`
		W          int    `json:"w"`
		H          int    `json:"h"`
	}

	// uploadのレスポンス
	uploadResponse struct {
		Media_id        int64  `json:"media_id"`
		Media_id_string string `json:"media_id_string"`
		Media_key       string `json:"media_key"`
		Size            int64  `json:"size"`
		Expires         int64  `json:"expires_after_secs"`
		Image           image  `json:"image"`
	}

	Twitter struct {
		client *http.Client
	}
)

// inits Twitter struct.
// confPath: full-path to twitter config file
func NewTwitter(confPath string) Twitter {

	b, err := os.ReadFile(confPath)

	if err != nil {
		fmt.Println(err)
	}

	cred := &keys{}
	json.Unmarshal(b, cred)

	config := oauth1.NewConfig(cred.API_KEY, cred.API_SECRET)
	token := oauth1.NewToken(cred.ACCESS_TOKEN, cred.ACCESS_SECRET)
	client := config.Client(oauth1.NoContext, token)

	return Twitter{
		client: client,
	}
}

// upload media.
// paths : full-paths of images to upload
// returns : list of uploaded media_ids
func (t Twitter) upload(paths ...string) []string {
	var ids []string
	for _, path := range paths {
		// 画像ファイルを開く
		file, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		// extract file name from full-path.
		filename := filepath.Base(path)

		// init io.Writer interface
		body := &bytes.Buffer{}

		// init multipart.Writer interface
		writer := multipart.NewWriter(body)

		// create form-data header
		part, err := writer.CreateFormFile("media", filename)
		if err != nil {
			fmt.Println(err)
		}

		// copy file content to `part` io.Writer.
		_, err = io.Copy(part, file)
		if err != nil {
			fmt.Println(err)
		}

		// add field data
		err = writer.WriteField("media_category", "tweet_image")
		if err != nil {
			fmt.Println(err)
		}

		// closing multipart.Writer writes trailing boundary.
		// So you need to close it here.
		err = writer.Close()

		if err != nil {
			fmt.Println(err)
		}

		// create POST request.
		req, err := http.NewRequest("POST", URL_IMAGE, body)
		if err != nil {
			fmt.Println(err)
		}

		// set Content-Type with Writers boundary.
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// execute request
		res, err := t.client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()

		if !(res.StatusCode >= 200 && res.StatusCode <= 299) {
			fmt.Printf("status-code:%v skipping...\n", res.StatusCode)
			continue
		}

		// read response body and log it to console.
		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}

		// parse resBody to `uploadResponse`
		resData := &uploadResponse{}
		json.Unmarshal(resBody, resData)

		// add media_id_string to ids
		ids = append(ids, resData.Media_id_string)
	}
	return ids
}

// msg   : message to tweet.
// paths : full-paths of images to upload
func (t Twitter) Tweet(msg string, paths ...string) {

	var ids []string // upload media_ids

	// 画像パスありならアップロード
	if len(paths) > 0 {
		ids = t.upload(paths...)
	}

	// body部 テキスト部分を設定
	data := map[string]interface{}{"text": msg}

	// upload画像がある場合は設定
	if ids != nil {
		data["media"] = map[string][]string{"media_ids": ids}
	}

	// bodyをJSON文字列に変換
	reqbody, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	// post requestを生成しヘッダ追加
	req, err := http.NewRequest("POST", URL_TWEET, strings.NewReader(string(reqbody)))
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/json") // 指定しないとエラーになる。

	// リクエスト実行
	res, err := t.client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}
