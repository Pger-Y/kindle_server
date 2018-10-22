package command

import (
	"github.com/kindle_server/types"
	"github.com/kindle_server/worker"
	"log"
	"strings"
)

type Command struct {
	keywords map[string]struct{}
	worker   worker.Worker
	split    string
}

func New(s string, w worker.Worker) *Command {
	c := &Command{
		keywords: map[string]struct{}{
			"register": struct{}{},
			"send":     struct{}{},
			"search":   struct{}{},
			"help":     struct{}{},
		},
		split:  s,
		worker: w,
	}
	return c
}

/*
func (c *Command) Employ(w worker.Worker) {
	c.worker = w
}
*/

//func (c *Command) Process(userid, message_type, value string) string {
func (c *Command) Process(vx_req *types.Xml) *types.Xml {
	var vx_resp types.Xml
	vx_resp.ToUserName, vx_resp.FromUserName = vx_req.FromUserName, vx_req.ToUserName
	vx_resp.MsgType = "text"
	vx_resp.CreateTime = vx_req.CreateTime + 1
	if vx_req.MsgType == "event" && vx_req.Event == "subscribe" {
		vx_resp.Content = c.worker.Usage()
	} else {
		values := strings.Split(vx_req.Content, c.split)
		ret, err := c.worker.Exec(vx_req.FromUserName, vx_req.MsgType, values...)
		vx_resp.Content = ret
		if err != nil {
			log.Println("Error while process ", vx_req, err)
		}
	}
	return &vx_resp

}
