package officialaccount

import (
	"encoding/xml"
	"time"
)

type Context struct {
	Message ReceiveMessage
}

type MsgBase struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	MsgType      string `xml:"MsgType"`
}

type MsgNews struct {
	MsgBase
	ArticleCount int           `xml:"ArticleCount"`
	Articles     []MsgNewsItem `xml:"Articles"`
}

type MsgNewsItem struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"Title"`
	Description string   `xml:"Description"`
	PicURL      string   `xml:"PicUrl"`
	URL         string   `xml:"Url"`
}

func (c *Context) ReplyNone(exit ...bool) *ReplyMessage {
	e := false
	if len(exit) > 0 {
		e = exit[0]
	}

	return &ReplyMessage{
		Reply: false,
		Exit:  e,
	}
}

func (c *Context) ReplyText(content string) *ReplyMessage {
	return c.appendData(&ReplyMessage{
		Reply: true,
		Exit:  true,
		Data: map[string]interface{}{
			"MsgType": "text",
			"Content": content,
		},
	})
}

func (c *Context) ReplyImage(mediaID string) *ReplyMessage {
	return c.appendData(&ReplyMessage{
		Reply: true,
		Exit:  true,
		Data: map[string]interface{}{
			"MsgType": "image",
			"MediaId": mediaID,
		},
	})
}

func (c *Context) ReplyVoice(mediaID string) *ReplyMessage {
	return c.appendData(&ReplyMessage{
		Reply: true,
		Exit:  true,
		Data: map[string]interface{}{
			"MsgType": "voice",
			"MediaId": mediaID,
		},
	})
}

func (c *Context) ReplyVideo(mediaID, title, description string) *ReplyMessage {
	return c.appendData(&ReplyMessage{
		Reply: true,
		Exit:  true,
		Data: map[string]interface{}{
			"MsgType":     "video",
			"MediaId":     mediaID,
			"Title":       title,
			"Description": description,
		},
	})
}

func (c *Context) ReplyMusic(title, description, musicURL, HQMusicURL, thumbMediaID string) *ReplyMessage {
	return c.appendData(&ReplyMessage{
		Reply: true,
		Exit:  true,
		Data: map[string]interface{}{
			"MsgType":      "music",
			"Title":        title,
			"Description":  description,
			"MusicURL":     musicURL,
			"HQMusicUrl":   HQMusicURL,
			"ThumbMediaId": thumbMediaID,
		},
	})
}

func (c *Context) ReplyNews(news []map[string]interface{}) *ReplyMessage {
	var data []interface{}
	for _, v := range news {
		data = append(data, v)
	}
	return c.appendData(&ReplyMessage{
		Reply: true,
		Exit:  true,
		Data: map[string]interface{}{
			"MsgType":      "news",
			"ArticleCount": len(news),
			"Articles":     map[string]interface{}{"item": data},
		},
	})
}

func (c *Context) appendData(message *ReplyMessage) *ReplyMessage {
	message.Data["ToUserName"] = c.Message["FromUserName"]
	message.Data["FromUserName"] = c.Message["ToUserName"]
	message.Data["CreateTime"] = time.Now().Unix()
	return message
}
