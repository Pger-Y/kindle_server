package command

import (
	"fmt"
	"log"
	"strings"

	"github.com/kindle_server/types"
	"github.com/kindle_server/worker"
)

type Command struct {
	keywords    map[string]struct{}
	c2w         map[string]worker.Worker
	split       string
	defaultFunc func(uid uint64, info string) (string, error)
}

func New(s string) *Command {
	c := &Command{
		keywords: map[string]struct{}{
			"register": struct{}{},
			"send":     struct{}{},
			"search":   struct{}{},
			"help":     struct{}{},
		},
		split:       s,
		c2w:         map[string]worker.Worker{},
		defaultFunc: nil,
	}
	return c
}

/*
func (c *Command) Employ(w worker.Worker) {
	c.worker = w
}
*/

//func (c *Command) Process(userid, message_type, value string) string {
func (c *Command) Usage() string {
	clist := []string{}
	for k, w := range c.c2w {
		u := w.Usage()
		m := fmt.Sprintf("%s命令的用法如下:%s", k, u)
		clist = append(clist, m)
	}
	m := strings.Join(clist, "\n")
	return m
}

func (c *Command) AddWorker(key string, w worker.Worker) {
	c.c2w[key] = w
}

func (c *Command) SetDefault(f func(uint64, string) (string, error)) {
	c.defaultFunc = f
}
func (c *Command) Process(vx_req *types.Xml) *types.Xml {
	var vx_resp types.Xml
	vx_resp.ToUserName, vx_resp.FromUserName = vx_req.FromUserName, vx_req.ToUserName
	vx_resp.MsgType = "text"
	vx_resp.CreateTime = vx_req.CreateTime + 1
	uid_str := types.Uid(vx_req.FromUserName)
	uid := uid_str.Hash()
	if vx_req.MsgType == "event" && vx_req.Event == "subscribe" {
		vx_resp.Content = c.Usage()
	} else {

		values := strings.Split(vx_req.Content, c.split)
		worker, ok := c.c2w[values[0]]
		if !ok || len(values) < 2 {
			if c.defaultFunc != nil {
				message, err := c.defaultFunc(uid, vx_req.Content)
				if err == nil {
					vx_resp.Content = message
				} else {
					vx_resp.Content = "默认功能无法进行处理，请进行检查后重试"
				}
			} else {
				//TODO
				// 增加相关功能的displa
				vx_resp.Content = "未找到相关功能，检查输入以及格式"
			}
		} else {
			ret, err := worker.Exec(vx_req.FromUserName, vx_req.MsgType, values[1:]...)
			vx_resp.Content = ret
			if err != nil {
				log.Println("Error while process ", vx_req, err)
			}
		}
	}
	return &vx_resp

}
