package internal

import (
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
		if v != nil {
			room := v.(*Room)
			p.ExitFromRoom(room)
			log.Debug("Agent删除房间资源~")
		}
	}
}
