package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/kindle_server/command"
	"github.com/kindle_server/config"
	"github.com/kindle_server/store"
	"github.com/kindle_server/types"
	"github.com/kindle_server/worker/account"
	"github.com/kindle_server/worker/kindle"
	"github.com/kindle_server/worker/kindle/mem"
	"log"
	"net/http"
)

func main() {
	r := gin.Default()
	ctx := context.Background()
	cfg, _, err := config.LoadFile("config.yaml")
	if err != nil {
		log.Fatal("Load Error:", err)
	}
	users := mem.NewUsers(ctx, cfg)
	store, err := store.NewStore(cfg.MySQL)
	if err != nil {
		log.Fatal("Error while connect mysql", err)
	}
	work_kindle := kindle.NewKindleWorker(users, store)
	account := account.NewAccount()
	cmd := command.New(cfg.Split)
	cmd.AddWorker("kindle", work_kindle)
	cmd.AddWorker("cal", account)

	r.POST("Xweixin_pathX", func(c *gin.Context) {
		var vx_req types.Xml
		if err := c.ShouldBindXML(&vx_req); err != nil {
			log.Println("Parse request error", err)

		} else {
			vx_resp := cmd.Process(&vx_req)
			c.XML(http.StatusOK, *vx_resp)
		}
	})
	r.Run()
}
