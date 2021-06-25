package internal

import (
	"fctbj/msg"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
)

func init() {
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)
}

func rpcNewAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	log.Debug("<-------------新链接请求连接--------------->")

	p := &Player{}
	p.Init()
	p.ConnAgent = a
	p.ConnAgent.SetUserData(p)
}

func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	p, ok := a.UserData().(*Player)
	if ok && p.ConnAgent == a {
		log.Debug("<-------------%v 主动断开链接--------------->", p.Id)
		rid, _ := hall.UserRoom.Load(p.Id)
		v, _ := hall.RoomRecord.Load(rid)
		if v == nil {
			room := v.(*Room)
			hall.UserRecord.Delete(p.Id)
			hall.RoomRecord.Delete(room.RoomId)
			hall.UserRoom.Delete(p.Id)
			c2c.UserLogoutCenter(p.Id, p.Password, p.Token) //todo
			leaveHall := &msg.Logout_S2C{}
			p.SendMsg(leaveHall)
			p.ConnAgent.Close()
		}
	}
}
