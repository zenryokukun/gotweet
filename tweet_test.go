package gotweet

import "testing"

func TestTweet(t *testing.T) {
	twitter := NewTwitter("./twitter_conf.json")
	twitter.Tweet("hello,world")
}

func TestUpload(t *testing.T) {
	twitter := NewTwitter("./twitter_conf.json")
	twitter.Tweet("hello,world", "./view.JPG", "./hill.JPG")
}
