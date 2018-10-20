package main

import (
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Xml struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int
	MsgType      string
	Content      string
	Event        string
}

func main() {
	r := gin.Default()
	r.POST("Xweixin_pathX", func(c *gin.Context) {
		var vx_req Xml
		if err := c.ShouldBindXML(&vx_req); err != nil {
			//c.XML(http.StatusOK,
			log.Println("Parse request error", err)

		}
		log.Println("Request:", vx_req)
		var vx_resp Xml
		vx_resp.ToUserName, vx_resp.FromUserName = vx_req.FromUserName, vx_req.ToUserName
		vx_resp.MsgType = "text"
		vx_resp.CreateTime = vx_req.CreateTime + 1
		vx_resp.Content = "We received you message:" + vx_req.Content
		s, err := xml.Marshal(vx_resp)
		log.Println("xml:marshal", err)
		log.Println(string(s))
		c.XML(http.StatusOK, vx_resp)
	})
	r.Run()
}
