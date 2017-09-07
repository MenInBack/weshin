package message

type Message struct {
	Meta
	Content
}

type Meta struct {
	FromUserName string `xml:"FromUserName,cdata"`
	ToUserName   string `xml:"ToUserName,cdata"`
	CreateTime   int64  `xml:"CreateTime,cdata"`
	MessageType  string `xml:"MsgType,cdata"`
}

type Content interface {
	GetMessageID() int64
}

type Text struct {
	Content   string `xml:"Content,cdata"`
	MessageID int64  `xml:"MsgId,cdata"`
}

func (c *Text) GetMessageID() int64 {
	if c == nil {
		return 0
	}
	return c.MessageID
}
