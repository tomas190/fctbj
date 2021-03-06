package internal

import (
	"fctbj/msg"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"reflect"
	"time"
)

func init() {
	handlerReg(&msg.Ping{}, handlePing)

	handlerReg(&msg.Login_C2S{}, handleLogin)
	handlerReg(&msg.Logout_C2S{}, handleLogout)

	handlerReg(&msg.JoinRoom_C2S{}, handleJoinRoom)

	handlerReg(&msg.PlayerAction_C2S{}, handlePlayerAction)
	handlerReg(&msg.ActionResult_C2S{}, handleActionResult)
	handlerReg(&msg.ProgressBar_C2S{}, handleProgressBar)
	handlerReg(&msg.LuckyPig_C2S{}, handleLuckyPig)
	handlerReg(&msg.PickUpGold_C2S{}, handlePickUpGold)

	handlerReg(&msg.ChangeRoomCfg_C2S{}, handleChangeRoomCfg)
	handlerReg(&msg.SendCoordinate_C2S{}, handleSendCoordinate)

}

// 注册消息处理函数
func handlerReg(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func handlePing(args []interface{}) {
	a := args[1].(gate.Agent)

	pingTime := time.Now().UnixNano() / 1e6
	pong := &msg.Pong{
		ServerTime: pingTime,
	}
	a.WriteMsg(pong)
}

func handleLogin(args []interface{}) {
	m := args[0].(*msg.Login_C2S)
	a := args[1].(gate.Agent)

	log.Debug("handleLogin 用户登入游戏~ :%v", m.Id)
	v, ok := hall.UserRecord.Load(m.Id)
	if ok { // 说明用户已存在
		p := v.(*Player)
		if p.ConnAgent == a { // 用户和链接都相同
			log.Debug("同一用户相同连接重复登录~")
			return
		} else { // 用户相同，链接不相同
			err := hall.ReplacePlayerAgent(p.Id, a)
			if err != nil {
				log.Error("用户链接替换错误", err)
			}

			c2c.UserLoginCenter(m.GetId(), m.GetPassWord(), m.GetToken(), func(au *Player) { // todo
				user, _ := hall.UserRecord.Load(p.Id)
				if user != nil {
					u := user.(*Player)
					u.NickName = au.NickName
					u.HeadImg = au.HeadImg
					u.Account = au.Account
					log.Debug("读取玩家历史数据:%v", u)
					login := &msg.Login_S2C{}
					rid, _ := hall.UserRoom.Load(p.Id)
					v, _ := hall.RoomRecord.Load(rid)
					if v != nil {
						login.IsBack = true
					}
					login.PlayerInfo = new(msg.PlayerInfo)
					login.PlayerInfo.Id = u.Id
					login.PlayerInfo.NickName = u.NickName
					login.PlayerInfo.HeadImg = u.HeadImg
					login.PlayerInfo.Account = u.Account
					a.WriteMsg(login)

					p.ConnAgent = a
					p.ConnAgent.SetUserData(u)

					p.OffLineTime = -1
					p.Password = m.GetPassWord()
					p.Token = m.GetToken()

					hall.OnlineUser.Store(u.Id, u)
				}
				// 判断玩家是否返回房间数据
				p.RespEnterRoom()
			})
		}
	} else if !hall.agentExist(a) { // 玩家首次登入
		c2c.UserLoginCenter(m.GetId(), m.GetPassWord(), m.GetToken(), func(u *Player) { //todo
			log.Debug("玩家首次登陆:%v", u)
			login := &msg.Login_S2C{}
			login.IsBack = false
			login.PlayerInfo = new(msg.PlayerInfo)
			login.PlayerInfo.Id = u.Id
			login.PlayerInfo.NickName = u.NickName
			login.PlayerInfo.HeadImg = u.HeadImg
			login.PlayerInfo.Account = u.Account
			a.WriteMsg(login)

			u.Init()
			// 重新绑定信息
			u.ConnAgent = a
			a.SetUserData(u)

			u.OffLineTime = -1
			u.Password = m.GetPassWord()
			u.Token = m.GetToken()

			hall.UserRecord.Store(u.Id, u)
			hall.OnlineUser.Store(u.Id, u)
		})
	}
}

func handleLogout(args []interface{}) {
	a := args[1].(gate.Agent)

	p, ok := a.UserData().(*Player)
	log.Debug("handleLeaveHall 玩家退出大厅~ : %v", p.Id)
	if ok {
		rid, _ := hall.UserRoom.Load(p.Id)
		v, _ := hall.RoomRecord.Load(rid)
		if v != nil {
			room := v.(*Room)
			p.ExitFromRoom(room)
			log.Debug("Logout删除房间资源~")
		} else {
			log.Debug("离开房间失败, 当前玩家未在房间内~")
		}
	}
}

func handleJoinRoom(args []interface{}) {
	m := args[0].(*msg.JoinRoom_C2S)
	a := args[1].(gate.Agent)

	p, ok := a.UserData().(*Player)
	log.Debug("handleJoinRoom 玩家加入房间~ : %v", p.Id)

	if ok {
		p.PlayerJoinRoom(m.Cfg)
	}
}

func handlePlayerAction(args []interface{}) {
	m := args[0].(*msg.PlayerAction_C2S)
	a := args[1].(gate.Agent)

	p, ok := a.UserData().(*Player)
	log.Debug("handlePlayerAction 玩家开始行动~ : %v", p.Id)

	if ok {
		p.PlayerAction(m)
	}
}

func handleActionResult(args []interface{}) {
	m := args[0].(*msg.ActionResult_C2S)
	a := args[1].(gate.Agent)

	p, ok := a.UserData().(*Player)
	log.Debug("handleActionResult 玩家行动结算~ : %v", p.Id)

	if ok {
		p.PlayerResult(m)
	}
}

func handleProgressBar(args []interface{}) {
	m := args[0].(*msg.ProgressBar_C2S)
	a := args[1].(gate.Agent)

	p, ok := a.UserData().(*Player)
	log.Debug("handleProgressBar 获取进度条金币~ : %v", p.Id)

	if ok {
		p.ProgressBetResp(m)
	}
}

func handleLuckyPig(args []interface{}) {
	a := args[1].(gate.Agent)

	p, ok := a.UserData().(*Player)
	log.Debug("handleLuckyPig 获取财运满满~ : %v", p.Id)

	if ok {
		p.WinLuckyPig()
	}
}

func handlePickUpGold(args []interface{}) {
	m := args[0].(*msg.PickUpGold_C2S)
	a := args[1].(gate.Agent)

	p, ok := a.UserData().(*Player)
	log.Debug("handleProgressBar 获取财神接金币~ : %v", p.Id)

	if ok {
		p.GodPickUpGold(m.BetNum)
	}
}

func handleChangeRoomCfg(args []interface{}) {
	m := args[0].(*msg.ChangeRoomCfg_C2S)
	a := args[1].(gate.Agent)

	p, ok := a.UserData().(*Player)
	log.Debug("handleChangeRoomCfg 修改区分配置~ : %v", p.Id)

	if ok {
		p.ChangeRoomCfg(m)
	}
}

func handleSendCoordinate(args []interface{}) {
	m := args[0].(*msg.SendCoordinate_C2S)
	a := args[1].(gate.Agent)

	p, ok := a.UserData().(*Player)
	log.Debug("handleSendCoordinate 保存位置信息~ : %v", p.Id)

	if ok {
		p.SaveCoordinate(m)
	}
}
