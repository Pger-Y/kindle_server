package store

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

type MySQLC struct {
	conn *sql.DB
}

func NewMySQL(h, u, p, port, db string) (*MySQLC, error) {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", u, p, h, port, db)
	c, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}
	m := &MySQLC{
		conn: c,
	}
	return m, nil
}

func (m *MySQLC) Replace(tb string, col []string, values ...interface{}) error {
	c_str := strings.Join(col, ",")
	p_list := make([]string, len(col))
	for i, _ := range col {
		p_list[i] = "?"
	}
	p_str := strings.Join(p_list, ",")
	sql_str := fmt.Sprintf("replace into %s(%s) values(%s)", tb, c_str, p_str)
	_, err := m.conn.Exec(sql_str, values...)
	return err
}

func (m *MySQLC) Select(sql_str string, f func(err error, values ...interface{}) error, v ...interface{}) error {
	rows, err := m.conn.Query(sql_str)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(v...); err != nil {
			return err
		}
		if err1 := f(err, v...); err1 != nil {
			return err1
		}
	}

	return nil
}

func (m *MySQLC) SelectOne(sql_str string, values ...interface{}) error {
	err := m.conn.QueryRow(sql_str).Scan(values...)
	return err
}
