package internal

import (
	"github.com/name5566/leaf/gate"
)

type Player struct {
	// 玩家代理链接
	ConnAgent gate.Agent

	Id        string
	NickName  string
	HeadImg   string
	Account   float64
	Password  string
	Token     string
	RoundId   string
	PackageId uint16

	WinResultMoney  float64
	LoseResultMoney float64
}

func (p *Player) Init() {
	p.RoundId = ""

}

//SendMsg 玩家向客户端发送消息
func (p *Player) SendMsg(msg interface{}) {
	if p.ConnAgent != nil {
		p.ConnAgent.WriteMsg(msg)
	}
}

