package types

import (
	"encoding/xml"
)

type Xml struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int
	Event        string
	MsgType      string
	Content      string
}
