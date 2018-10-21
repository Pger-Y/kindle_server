package mem

import (
	"context"
	"github.com/kindle_server/config"
	"github.com/kindle_server/types"
	"sync"
	"time"
)

type Users struct {
	uinfos     map[uint64]*types.UserInfo
	mtx        sync.RWMutex
	gcInterval types.Duration // default 1m
	lifetime   types.Duration //default 30m
}

func NewUsers(ctx context.Context, config *config.GlobalConfig) *Users {
	us := &Users{
		uinfos:     map[uint64]*types.UserInfo{},
		gcInterval: config.GcInterval,
		lifetime:   config.LifeTime,
	}
	us.Run(ctx)
	return us
}

func (us *Users) Put(u *types.UserInfo) {
	us.mtx.Lock()
	defer us.mtx.Unlock()
	k := u.Hash()
	us.uinfos[k] = u
}

func (us *Users) Get(uid string) (*types.UserInfo, bool) {
	us.mtx.RLock()
	defer us.mtx.RUnlock()
	t := types.Uid(uid)
	u := t.Hash()
	user, ok := us.uinfos[u]
	return user, ok
	//false means no register yet

}

func (us *Users) Run(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(time.Duration(us.gcInterval))
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				us.gc()
			}
		}
	}()
}

func (us *Users) gc() {
	us.mtx.Lock()
	defer us.mtx.Unlock()
	now := time.Now()
	for uid, user := range us.uinfos {
		if user.Atime.Add(time.Duration(us.lifetime)).Before(now) {
			delete(us.uinfos, uid)
		}
	}
}
