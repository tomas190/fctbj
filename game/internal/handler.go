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
	handlerReg(&msg.SendWinMoney_C2S{}, handleSendWinMoney)
	handlerReg(&msg.GetRewards_C2S{}, handleGetRewards)

	handlerReg(&msg.ChangeRoomCfg_S2C{}, handleChangeRoomCfg)
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
			//ErrorResp(a, msg.ErrorMsg_UserRepeatLogin, "重复登录")
			return
		} else { // 用户相同，链接不相同
			err := hall.ReplacePlayerAgent(p.Id, a)
			if err != nil {
				log.Error("用户链接替换错误", err)
			}

			login := &msg.Login_S2C{}
			user, _ := hall.UserRecord.Load(p.Id)
			if user != nil {
				u := user.(*Player)
				login.PlayerInfo = new(msg.PlayerInfo)
				login.PlayerInfo.Id = u.Id
				login.PlayerInfo.NickName = u.NickName
				login.PlayerInfo.HeadImg = u.HeadImg
				login.PlayerInfo.Account = u.Account
				a.WriteMsg(login)

				p.ConnAgent = a
				p.ConnAgent.SetUserData(u)
			}

			rid, _ := hall.UserRoom.Load(p.Id)
			v, _ := hall.RoomRecord.Load(rid)
			if v != nil {
				room := v.(*Room)
				roomData := room.RespRoomData()
				enter := &msg.EnterRoom_S2C{}
				enter.RoomData = roomData
				p.SendMsg(enter)
				log.Debug("返回当前房间~")
			}
		}
	} else if !hall.agentExist(a) { // 玩家首次登入
		//u := &Player{}
		//u.Id = m.Id
		//u.Account = 10000
		//u.NickName = m.Id
		//u.Token = m.Token
		//
		//login := &msg.Login_S2C{}
		//login.PlayerInfo = new(msg.PlayerInfo)
		//login.PlayerInfo.Id = u.Id
		//login.PlayerInfo.NickName = u.NickName
		//login.PlayerInfo.HeadImg = u.HeadImg
		//login.PlayerInfo.Account = u.Account
		//a.WriteMsg(login)
		//
		//u.Init()
		//// 重新绑定信息
		//u.ConnAgent = a
		//a.SetUserData(u)
		//
		//u.Password = m.GetPassWord()
		//u.Token = m.GetToken()
		//
		//hall.UserRecord.Store(u.Id, u)

		c2c.UserLoginCenter(m.GetId(), m.GetPassWord(), m.GetToken(), func(u *Player) { //todo

			log.Debug("玩家首次登陆:%v", u.Id)
			login := &msg.Login_S2C{}
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

			u.Password = m.GetPassWord()
			u.Token = m.GetToken()

			hall.UserRecord.Store(u.Id, u)
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
			hall.UserRecord.Delete(p.Id)
			log.Debug("删除房间资源~")
			hall.RoomRecord.Delete(room.RoomId)
			hall.UserRoom.Delete(p.Id)
			c2c.UserLogoutCenter(p.Id, p.Password, p.Token) //todo
			leaveHall := &msg.Logout_S2C{}
			p.SendMsg(leaveHall)
			p.ConnAgent.Close()
		} else {
			log.Debug("离开房间失败, 当前玩家未在房间内~")
		}
	}
}

func handleJoinRoom(args []interface{}) {
	a := args[1].(gate.Agent)

	p, ok := a.UserData().(*Player)
	log.Debug("handleJoinRoom 玩家加入房间~ : %v", p.Id)

	if ok {
		p.PlayerJoinRoom()
	}
}

func handlePlayerAction(args []interface{}) {
	m := args[0].(*msg.PlayerAction_C2S)
	a := args[1].(gate.Agent)

	p, ok := a.UserData().(*Player)
	log.Debug("handlePlayerAction 玩家开始行动~ : %v", p.Id)

	if ok {
		log.Debug("玩家行動:%v", m)
		p.PlayerAction(m.DownBet)
	}
}

func handleSendWinMoney(args []interface{}) {
	m := args[0].(*msg.SendWinMoney_C2S)
	a := args[1].(gate.Agent)

	p, ok := a.UserData().(*Player)
	log.Debug("handleSendWinMoney 玩家发送赢钱~ : %v", p.Id)

	if ok {
		p.GetPlayerWinMoney(m.Money)
	}
}

func handleGetRewards(args []interface{}) {
	a := args[1].(gate.Agent)

	p, ok := a.UserData().(*Player)
	log.Debug("handleGetRewards 玩家获取奖源~ : %v", p.Id)

	if ok {
		p.GetRewardsInfo()
	}
}

func handleChangeRoomCfg(args []interface{}) {
	m := args[0].(*msg.ChangeRoomCfg_C2S)
	a := args[1].(gate.Agent)

	p, ok := a.UserData().(*Player)
	log.Debug("handleChangeRoomCfg 修改区分配置~ : %v", p.Id)

	if ok {
		rid, _ := hall.UserRoom.Load(p.Id)
		v, _ := hall.RoomRecord.Load(rid)
		if v != nil {
			room := v.(*Room)
			room.Config = m.Config

			// 发送配置数据
			data := &msg.ChangeRoomCfg_S2C{}
			p.SendMsg(data)
		}
	}
}
