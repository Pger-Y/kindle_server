package store

import (
	"github.com/kindle_server/config"
	"github.com/kindle_server/types"
)

type Store struct {
	tb_user  string
	col_user []string
	cli      *MySQLC
}

func NewStore(c *config.MySQLConfig) (*Store, error) {
	cli, err := NewMySQL(c.Host, c.User, c.Password, c.Port, c.Database)
	if err != nil {
		return nil, err
	}
	s := &Store{
		tb_user:  "user_info",
		col_user: []string{"userid", "kindle_address", "mail_address", "mail_passwd", "smtp_server"},
		cli:      cli,
	}
	return s, nil

}

func (s *Store) User2Sql(u *types.UserInfo) error {
	err := s.cli.Replace(s.tb_user, s.col_user, u.Userid, u.KindleAddress)
	return err
}
