package internal

import (
	"errors"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"sync"
	"time"
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

// 记录玩家数据
func (hall *GameHall) RecordPlayerData() {
	ticker := time.NewTicker(180 * time.Second)
	defer ticker.Stop()

	for { // 循环每3秒处理玩家数据
		select {
		case <-ticker.C:
			hall.UserRecord.Range(func(key, value interface{}) bool {
				v := value.(*Player)
				// 判断玩家期间是否行动过
				if v.DownBet > 0 {
					// 插入玩家数据
					v.HandlePlayerData()
				}
				return true
			})
		}
	}
}

// 处理房间数据
func (hall *GameHall) HandleRoomData() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tn := time.Now().Hour()
			hall.UserRecord.Range(func(key, value interface{}) bool {
				v := value.(*Player)
				if v.OffLineTime != -1 {
					if v.OffLineTime-tn == 0 {
						hall.UserRecord.Delete(v.Id)
						hall.UserRoom.Delete(v.Id)
						rid, _ := hall.UserRoom.Load(v.Id)
						v, _ := hall.RoomRecord.Load(rid)
						if v != nil {
							room := v.(*Room)
							hall.RoomRecord.Delete(room.RoomId)
						}
					}
				}
				return true
			})
		}
	}
}
