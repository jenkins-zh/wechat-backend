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
	MsgType string `json:"msgType"`
	Kind    string `json:"kind"`
}

type TextResponseBody struct {
	ResponseBody
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	Content      string
}

type NewsResponseBody struct {
	ResponseBody
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	ArticleCount int
	Articles     Articles
}

type Articles struct {
	XMLName  xml.Name  `xml:"Articles"`
	Articles []Article `xml:"item"`
}

type Article struct {
	Title       string
	Description string
	PicUrl      string
	Url         string
}
