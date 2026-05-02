package u

import (
	"encoding/xml"
	"time"
)

// Baseer 1. 定义通用接口，作为「父类类型」
type Baseer interface {
	GetMsgType() string
	GetToUserName() string
	GetFromUserName() string
	SetTarget(target string)
}

type BaseMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `json:"toUserName,omitempty" xml:"ToUserName"`     // 开发者微信号（公众号的原始ID）
	FromUserName string   `json:"fromUserName,omitempty" xml:"FromUserName"` // 发送方帐号（用户的OpenID）
	CreateTime   int64    `json:"createTime,omitempty" xml:"CreateTime"`     // 消息创建时间（时间戳）
	MsgType      string   `json:"msgType,omitempty" xml:"MsgType"`           // 消息类型（text, image, voice, video, shortvideo, location, link, event等）
	Target       string   `json:"target,omitempty" xml:"Target"`
}

func (this *BaseMessage) SetTarget(target string) {
	this.Target = target
}
func (this *BaseMessage) GetToUserName() string {
	return this.ToUserName
}

func (this *BaseMessage) GetFromUserName() string {
	return this.FromUserName
}
func (this *BaseMessage) GetMsgType() string {
	return this.MsgType
}

// TextMessage 定义接收到的微信文本消息结构
type TextMessage struct {
	BaseMessage
	Content string `json:"content,omitempty" xml:"Content"` // 文本消息内容
	MsgId   int64  `json:"msgId,omitempty" xml:"MsgId"`     // 消息ID（64位整型）
}

// EventMessage 示例：关注/取消关注事件消息结构
type EventMessage struct {
	BaseMessage
	Event    string `json:"event,omitempty" xml:"Event"`       // 事件类型（如 subscribe-关注, unsubscribe-取消关注, CLICK-点击菜单）
	EventKey string `json:"eventKey,omitempty" xml:"EventKey"` // 事件KEY值
}

// <xml>
// <ToUserName><![CDATA[gh_5f173d9e52e4]]></ToUserName>
// <FromUserName><![CDATA[oIin3168TLKg1X8OU2xBBWLlMEdI]]></FromUserName>
// <CreateTime>1777706248</CreateTime>
// <MsgType><![CDATA[image]]></MsgType>
// <PicUrl><![CDATA[http://mmecoa.qpic.cn/sz_mmecoa_jpg/elicVUxeBpKNTsCibiaTpt9gdRWM63W0WEhBn5HY3DmwAIx2g1bRnaJq3LHHicyAZQrJgJicyRZcvQmWcUXnZibo5yn1S82CIyT5XGDd1dwuuzk8g/0]]></PicUrl>
// <MsgId>25450632777398658</MsgId>
// <MediaId><![CDATA[7bQqknQ3S8AEEfKeb585LTmkbLfXYWS1YulnPlZIU56BvihHmIGEZ0l97sNRR92l]]></MediaId>
// </xml>
type ImageMessage struct {
	BaseMessage
	PicUrl  string `json:"picUrl,omitempty" xml:"PicUrl"`
	MsgId   string `json:"msgId,omitempty" xml:"MsgId"`
	MediaId string `json:"mediaId,omitempty" xml:"MediaId"`
}

type LinkMessage struct {
	BaseMessage
	Title       string `json:"title,omitempty" xml:"Title"`
	Description string `json:"description,omitempty" xml:"Description"`
	Url         string `json:"url,omitempty" xml:"Url"`
	MsgId       int64  `json:"msgId,omitempty" xml:"MsgId"`
	MsgDataId   string `json:"msgDataId,omitempty" xml:"MsgDataId"`
	Idx         string `json:"idx,omitempty" xml:"Idx"`
}

// TextResponse 定义回复给微信服务器的文本消息结构
type TextResponse struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATA    `xml:"ToUserName"`   // 接收方帐号（用户的OpenID）
	FromUserName CDATA    `xml:"FromUserName"` // 发送方帐号（公众号的原始ID）
	CreateTime   int64    `xml:"CreateTime"`   // 消息创建时间（时间戳）
	MsgType      CDATA    `xml:"MsgType"`      // 消息类型（此处为 "text"）
	Content      CDATA    `xml:"Content"`      // 回复的文本内容
}

type LinkResponse struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATA    `xml:"ToUserName"`
	FromUserName CDATA    `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      CDATA    `xml:"MsgType"` // 此处为 "event"
	Title        CDATA    `xml:"Title"`
	Description  CDATA    `xml:"Description"`
	Url          CDATA    `xml:"Url"`
	MsgId        int64    `xml:"MsgId"`
	MsgDataId    int64    `xml:"MsgDataId"`
	Idx          int64    `xml:"Idx"`
}

// CDATA 处理XML CDATA标签（如果需要生成回复，包含CDATA时有用）
type CDATA struct {
	Value string `xml:",cdata"`
}

// CreateTextResponse 辅助函数，用于创建文本回复消息结构体
func CreateTextResponse(toUser, fromUser, content string) TextResponse {
	return TextResponse{
		ToUserName:   CDATA{Value: toUser},
		FromUserName: CDATA{Value: fromUser},
		CreateTime:   time.Now().Unix(),
		MsgType:      CDATA{Value: "text"},
		Content:      CDATA{Value: content},
	}
}
