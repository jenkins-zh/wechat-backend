package main

import (
	"encoding/xml"
	"time"
)

type TextRequestBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string
	MsgId        int
	Event        string
}

type ResponseBody struct {
	Keyword string `json:"keyword"`

	MsgType      string `json:"msgType" yaml:"msgType" xml:"MsgType"`
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
}

type TextResponseBody struct {
	ResponseBody `yaml:",inline"`
	XMLName      xml.Name `xml:"xml"`
	Content      string
}

type NewsResponseBody struct {
	ResponseBody `yaml:",inline"`
	XMLName      xml.Name `xml:"xml"`
	ArticleCount int      `json:"articleCount" yaml:"articleCount" xml:"ArticleCount"`
	Articles     Articles `yaml:",inline"`
}

type ImageResponseBody struct {
	ResponseBody `yaml:",inline"`
	XMLName      xml.Name `xml:"xml"`
	Image        Image
}

type Articles struct {
	// XMLName  xml.Name  `xml:"Articles"`
	Articles []Article `xml:"item"`
}

type Image struct {
	MediaID string `json:"mediaId" yaml:"mediaID" xml:"MediaId"`
}

type Article struct {
	Title       string
	Description string
	PicUrl      string `json:"picUrl" yaml:"picUrl" xml:"PicUrl"`
	Url         string
}
