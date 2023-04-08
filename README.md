# GO-TWEET

Tweet用モジュール。Twitter API V2対応。

## これは何ですか？

GoでTweetするモジュールです。

対応しているエンドポイントは、2023.04.05時点で以下の通りです。
- ```https://api.twitter.com/2/tweets```

メディアアップロードも出来ますが、V2版はまだリリースされていません。公式案内のとおり、旧API(V1.1)を利用しています。
リリースされ次第対応します。

## いつ使いますか？

テキストのみのTweetしたいとき、画像付きでツイートしたいとき。

## どう使いますか？

### Twitterの認情報をjson形式でファイルに保存します。

```json
{
  "API_KEY": "MY-API-KEY",
  "API_SECRET": "MY-API-SECRET",
  "BEARER": "MY-BEARER",
  "ACCESS_TOKEN": "MY-ACCESS-TOKEN",
  "ACCESS_SECRET": "MY-ACCESS-SECRET"
}
```

### Goでimportして使います。

```go
package main

import "github.com/zenryokukun/gotweet"

func main(){
  // init gotweet with your credential file.
  twitter := gotweet.NewTwitter("path-to-credential-file")

  // post tweet with text.
  twitter.Tweet("Hello,World!")

  // post tweet with image(s).
  twitter.Tweet("Hello,World!","path-to-img1",..,"path-to-img4")
}

```
