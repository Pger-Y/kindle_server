package types

import (
	"fmt"
	"github.com/checkmail"
	"time"
)

const (
	offset64 = 14695981039346656037
	prime64  = 1099511628211
)

// hashNew initializies a new fnv64a hash value.
func hashNew() uint64 {
	return offset64
}

func hashAdd(h uint64, s string) uint64 {
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= prime64
	}
	return h
}

func NewUser(uid, kmail, email, passwd, ss string) *UserInfo {
	u := &UserInfo{
		Userid:        Uid(uid),
		KindleAddress: kmail,
		MailAddress:   email,
		Passwd:        passwd,
		SmtpServer:    ss,
		Atime:         time.Now(),
	}
	return u
}

type Uid string

func (ui *Uid) Hash() uint64 {
	sum := hashNew()
	sum = hashAdd(sum, string(*ui))
	return sum

}

type UserInfo struct {
	Userid        Uid
	KindleAddress string
	MailAddress   string
	Passwd        string
	SmtpServer    string
	Atime         time.Time
}

func (u *UserInfo) Hash() uint64 {
	return u.Userid.Hash()
}

func mail_validate(m string) error {
	if err := checkmail.ValidateFormat(m); err != nil {
		return err
	}
	// is too slow to varify host
	/*
		if err := checkmail.ValidateHost(m); err != nil {
			return err
		}
	*/
	return nil
}

func (u *UserInfo) Validate() error {
	if err := mail_validate(u.KindleAddress); err != nil {
		return fmt.Errorf("Error while validate kindle_mail:%v", u.KindleAddress)
	}
	if u.MailAddress != "" {
		if err := mail_validate(u.MailAddress); err != nil {
			return fmt.Errorf("Error while validate your mail address:%v", u.MailAddress)
		}
	}
	return nil
}
