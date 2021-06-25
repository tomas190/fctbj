package internal

import (
	"errors"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"sync"
)

type GameHall struct {
	UserRecord sync.Map // 用户记录
	RoomRecord sync.Map // 房间记录
	UserRoom   sync.Map // 用户房间
}

func NewHall() *GameHall {
	return &GameHall{
		UserRecord: sync.Map{},
		RoomRecord: sync.Map{},
		UserRoom:   sync.Map{},
	}
}


//ReplacePlayerAgent 替换用户链接
func (hall *GameHall) ReplacePlayerAgent(Id string, agent gate.Agent) error {
	log.Debug("用户重连或顶替，正在替换agent %+v", Id)
	// tip 这里会拷贝一份数据，需要替换的是记录中的，而非拷贝数据中的，还要注意替换连接之后要把数据绑定到新连接上
	if v, ok := hall.UserRecord.Load(Id); ok {
		//ErrorResp(agent, msg.ErrorMsg_UserRemoteLogin, "异地登录")
		user := v.(*Player)
		user.ConnAgent = agent
		user.ConnAgent.SetUserData(v)
		return nil
	} else {
		return errors.New("用户不在记录中~")
	}
}

//agentExist 链接是否已经存在
func (hall *GameHall) agentExist(a gate.Agent) bool {
	var exist bool
	hall.UserRecord.Range(func(key, value interface{}) bool {
		u := value.(*Player)
		if u.ConnAgent == a {
			exist = true
		}
		return true
	})
	return exist
}