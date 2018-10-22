package types

import (
	"log"
	"testing"
)

func TestUser(t *testing.T) {
	km := "reg_used@kindle.cn"
	em := "kindle_pusherx@163.com"
	passwd := "another85576909"
	ss := "smtp.163.com"
	u := NewUser("uidddd", km, em, passwd, ss)
	err := u.Validate()
	log.Println(err)
}
