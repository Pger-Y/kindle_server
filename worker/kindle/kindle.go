package kindle

import (
	"fmt"
	//"github.com/kindle_server/config"
	"github.com/kindle_server/store"
	"github.com/kindle_server/types"
	"github.com/kindle_server/worker/kindle/mem"
	//"log"
	"strings"
	//"time"
)

type KindleWorker struct {
	ucache *mem.Users
	store  *store.Store
}

func NewKindleWorker(uc *mem.Users, s *store.Store) *KindleWorker {
	k := &KindleWorker{
		ucache: uc,
		store:  s,
	}
	return k

}

func (kw *KindleWorker) Usage() string {
	s := fmt.Sprint("Usage: \n",
		"register|example@kindle.com|example@163.com|examplepassword|smtp.163.com register your information\n",
		"send(developing...)|www.example.com/example.mobi will download ebook and send it to your kindle mail\n",
		"search(developing...)|The great gatsby will search ebook in our database that others has downloaded,yes! you are improving our data quality\n",
		"feedback(developing...) will display the url you have used for download then type feedback|1233|score[1-10] to the resource",
	)
	return s
}

func (kw *KindleWorker) Send(uid string, infos ...string) (string, error) {
	return "", nil
}

func (kw *KindleWorker) Exec(uid string, msg_type string, infos ...string) (string, error) {
	var message_ret string
	var err error
	if len(infos) < 2 && strings.ToLower(infos[0]) != "help" {
		message_ret = "Parsing cmd error type help for usage!"
		err := fmt.Errorf("cmd too short")
		return message_ret, err
	} else {
		do := strings.ToLower(infos[0])
		switch do {
		case "help":
			message_ret = kw.Usage()
			err = nil
		case "register":
			message_ret, err = kw.Register(uid, infos[1:]...)
		case "send":
			message_ret, err = kw.Send(uid, infos[1:]...)
		default:
			message_ret = "Unsupport cmd,check manual again:\n" + kw.Usage()
			err = fmt.Errorf("Unsupport cmd %v", do)
		}
	}
	return message_ret, err
}

func (kw *KindleWorker) Register(uid string, infos ...string) (string, error) {
	var kmail, email, pwd, ss string
	var message string
	info_len := len(infos)
	if info_len < 4 {
		message = "Register info format error example: register:example@kindle.com|example@163.com|examplepasswd|smtp.163.com\n or [unsupport yet]register:example@kindle.com"
		err := fmt.Errorf("register is to short")
		return message, err
	}
	//TODO @just4fun.im and default auto complete suffix of 163 email
	kmail = infos[0]
	if info_len >= 3 {
		email = infos[1]
		pwd = infos[2]
	}
	if info_len > 3 {
		ss = infos[3]
	}
	if ss == "" {
		for _, dtag := range []string{"163.com", "126.com"} {
			if strings.Index(email, dtag) != -1 {
				ss = fmt.Sprintf("smtp.%s", dtag)
				break
			}
		}
	}
	u := types.NewUser(uid, kmail, email, pwd, ss)

	if err := u.Validate(); err != nil {
		message = "register info invalid,please double check"
		return message, err
	}
	kw.ucache.Put(u)
	err := kw.store.User2Sql(u)
	message = "Register Successfully"
	if err != nil {
		message = "Register failed please send a message"
	}
	return message, err
}
