package account

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	cal "github.com/Pger-Y/calculate-go"
	"github.com/kindle_server/types"
)

type Account struct {
	data map[uint64]map[string]float64
	mtx  sync.RWMutex
}

func NewAccount() *Account {
	a := &Account{
		data: map[uint64]map[string]float64{},
	}
	return a
}

func (a *Account) Usage() string {
	s := fmt.Sprintf("Usage:\n",
		"cal:save:last=0(设置初始值，用以存储后续计算结果)\n",
		"cal:expr:last-100+222(这种方式会将结果存储到last值中，后续运算可以继续使用)\n",
		"cal:100-20*3\n",
	)
	return s
}

func (a *Account) Expr(uid uint64, info string) (string, error) {
	a.mtx.RLock()
	defer a.mtx.RUnlock()
	var find_key string
	if m, ok := a.data[uid]; ok {
		for k, v := range m {
			if strings.Index(info, k) == -1 {
				continue
			}
			find_key = k
			v_s := strconv.FormatFloat(v, 'f', 2, 64)
			info = strings.Replace(info, k, v_s, -1)
		}
	}
	v, err := cal.Calculate(info)
	if err != nil {
		m := fmt.Sprintf("calculate error[%s],check expression", err.Error())
		return m, fmt.Errorf(m)
	} else {
		message := fmt.Sprintf("%s = %f", info, v)
		if find_key != "" {
			a.save(uid, find_key, v)
		}
		return message, nil
	}

}

func (a *Account) set(uid uint64, info string) (string, error) {
	args := strings.Split(info, "=")
	if len(args) < 2 {
		return types.ArgErrorMessage, types.ArgFormatError
	}
	var_value := args[0]
	v_str := args[1]
	if vf, err := strconv.ParseFloat(v_str, 64); err != nil {
		return types.ArgFormatErrorMessage, types.ArgFormatError
	} else {
		a.save(uid, var_value, vf)
	}

	return "设置成功", nil
}

func (a *Account) save(uid uint64, k string, v float64) {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	_, ok := a.data[uid]
	if !ok {
		a.data[uid] = map[string]float64{}
	}
	a.data[uid][k] = v
}

func (a *Account) Exec(uid, msg_type string, infos ...string) (string, error) {
	var message_ret string
	u := types.Uid(uid)
	uh := u.Hash()
	if len(infos) < 2 && strings.ToLower(infos[0]) != "help" {
		message_ret = "Parsing cmd error type help for usage!"
		err := fmt.Errorf("cmd too short")
		return message_ret, err
	} else {
		do := strings.ToLower(infos[0])
		switch do {
		case "save":
			return a.set(uh, infos[1])
		case "expr":
			return a.Expr(uh, infos[1])
		}
	}
	m := "Account Parse Error!"
	return m, fmt.Errorf(m)
}
