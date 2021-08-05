package internal

import (
	"fctbj/msg"
	"github.com/name5566/leaf/log"
)

func (p *Player) PlayerJoinRoom(cfgId string) {
	// 判断玩家信息是否为空
	if p.Id == "" {
		p.SendErrMsg(RECODE_PlayerInfoIDIsNull)
		return
	}

	rid, _ := hall.UserRoom.Load(p.Id)
	v, _ := hall.RoomRecord.Load(rid)
	if v != nil { // 当前玩家已经存在房间
		room := v.(*Room)
		enter := &msg.EnterRoom_S2C{}
		enter.RoomData = room.RespRoomData()
		// 判断该金币区间是否存在金币位置存储，如果存在则返回，不存在则返回空
		if p.ConfigPlace[room.Config] != nil {
			enter.IsChange = true
			enter.Coordinates = p.ConfigPlace[room.Config]
		} else {
			enter.IsChange = false
			enter.Coordinates = p.ConfigPlace[room.Config]
		}
		p.SendMsg(enter)
	} else {
		// 创建房间
		r := &Room{}
		r.Init()
		r.Config = cfgId
		hall.RoomRecord.Store(r.RoomId, r)
		hall.UserRoom.Store(p.Id, r.RoomId)

		p.RoomId = r.RoomId

		// 插入玩家信息 todo
		p.FindPlayerInfo()

		//返回前端房间信息
		data := &msg.JoinRoom_S2C{}
		data.RoomData = r.RespRoomData()
		p.SendMsg(data)
	}
}

func (p *Player) ExitFromRoom(room *Room) {
	// 判断玩家期间是否行动过
	if p.DownBet > 0 {
		// 插入玩家数据
		p.HandlePlayerData()
	}

	// 删除房间资源
	hall.UserRecord.Delete(p.Id)
	hall.RoomRecord.Delete(room.RoomId)
	hall.UserRoom.Delete(p.Id)
	c2c.UserLogoutCenter(p.Id, p.Password, p.Token) //todo
	leaveHall := &msg.Logout_S2C{}
	p.SendMsg(leaveHall)
	p.ConnAgent.Close()
}

func (p *Player) PlayerAction(downBet float64) {
	// 判断玩家信息是否为空
	if p.Id == "" {
		p.SendErrMsg(RECODE_PlayerInfoIDIsNull)
		log.Debug("玩家信息ID为空,不能进行下注")
		return
	}
	// 判断玩家金额是否足够
	if p.Account < downBet {
		p.SendErrMsg(RECODE_UserMoneyNotEnough)
		log.Debug("玩家金额不足,不能进行下注~")
		return
	}

	// 判断下注金币是否对应房间配置金额(防止刷钱)
	rid, _ := hall.UserRoom.Load(p.Id)
	v, _ := hall.RoomRecord.Load(rid)
	if v != nil {
		room := v.(*Room)
		if CfgMoney[room.Config] != downBet {
			log.Debug("房间配置金额:%v,%v", CfgMoney[room.Config], downBet)
			p.SendErrMsg(RECODE_RoomCfgMoneyERROR)
			log.Debug("房间配置金额不对!")
			return
		}
	}

	p.Account -= downBet
	p.DownBet += downBet
	p.TotalLoseMoney += downBet

	p.DownBetCount++
	var IsDown bool
	if p.DownBetCount >= 50 {
		IsDown = true
		p.DownBetCount = 0
	}

	// 先判断盈余池是否有钱，然后在处理玩家是否赢钱
	surMoney := GetSurPlusMoney()
	log.Debug("盈余池的金额:%v", surMoney)
	var goldNum int32
	var taxR float64
	//if surMoney >= 0 { //todo
	goldNum = p.randGoldNum()
	if goldNum > 0 {
		// 玩家赢钱结算
		winMoney := downBet * float64(goldNum)
		pac := packageTax[p.PackageId]
		taxR = pac / 100
		tax := winMoney * taxR
		resultMoney := winMoney - tax

		p.Account += resultMoney
		p.TotalWinMoney += winMoney

		log.Debug("获取赢钱的金额:%v", winMoney)
	}
	//}

	log.Debug("玩家当前winNum：%v", goldNum)
	data := &msg.PlayerAction_S2C{}
	data.WinNum = goldNum
	data.Account = p.Account
	data.Tax = taxR
	data.LuckyBag = IsDown
	p.SendMsg(data)
}

func (p *Player) randGoldNum() int32 {
	var goldNum int32
	num := RandInRange(0, 100)
	if num >= 0 && num <= 50 { //65
		goldNum = 0
	} else if num >= 51 && num <= 85 {
		goldNum = 1
	} else if num >= 86 && num <= 95 {
		goldNum = 2
	} else if num >= 96 && num <= 100 {
		goldNum = 3
	}
	return goldNum
}

func (p *Player) GetRewardsInfo() {
	rid, _ := hall.UserRoom.Load(p.Id)
	v, _ := hall.RoomRecord.Load(rid)
	if v != nil {
		room := v.(*Room)
		// 房间配置金额
		cfgMoney := CfgMoney[room.Config]

		data := &msg.GetRewards_S2C{}
		var winMoney float64
		num := RandInRange(0, 100)
		if num >= 0 && num <= 5 {
			data.RewardsNum = GOLD
		} else if num >= 6 && num <= 12 {
			data.RewardsNum = RICH
			data.GetMoney = GetRICH(cfgMoney)
			winMoney = data.GetMoney
		} else if num >= 13 && num <= 30 {
			data.RewardsNum = PUSH
			data.GetMoney = GetPUSH(cfgMoney)
			winMoney = data.GetMoney
		} else if num >= 31 && num <= 100 {
			data.RewardsNum = LUCKY
			data.LuckyPig = GetLUCKY(cfgMoney)
			winMoney = data.LuckyPig.PigSuccess
		}

		// 结算
		pac := packageTax[p.PackageId]
		taxR := pac / 100
		tax := winMoney * taxR
		resultMoney := winMoney - tax

		p.Account += resultMoney
		p.TotalWinMoney += winMoney

		// 发送小游戏获奖
		data.Account = p.Account
		p.SendMsg(data)

		log.Debug("获取赢钱的金额:%v", winMoney)
	}
}

func (p *Player) ProgressBetResp(bet int32) {
	p.ProgressBet += bet

	log.Debug("p.ProgressBet 长度为:%v", p.ProgressBet)

	var betNum int32
	if p.ProgressBet >= 3 && p.ProgressBet <= 5 {
		betNum = 1
		data := &msg.ProgressBar_S2C{}
		data.ProBar = betNum
		p.SendMsg(data)
	} else if p.ProgressBet >= 6 && p.ProgressBet <= 8 {
		betNum = 2
		data := &msg.ProgressBar_S2C{}
		data.ProBar = betNum
		p.SendMsg(data)
	} else if p.ProgressBet >= 10 {
		betNum = 6
		// 发送进度条
		data := &msg.ProgressBar_S2C{}
		data.ProBar = betNum
		p.SendMsg(data)
		// 小游戏执行
		p.GetRewardsInfo()
		p.ProgressBet = 0
	}

	//var betNum int32
	//rid, _ := hall.UserRoom.Load(p.Id)
	//v, _ := hall.RoomRecord.Load(rid)
	//if v != nil {
	//	room := v.(*Room)
	//	// 房间配置金额
	//	money := CfgMoney[room.Config]
	//	surMoney := GetSurPlusMoney()
	//
	//	// 盈余池金额足够小游戏获奖时
	//	log.Debug("获奖的估计金额:%v,盈余池金额:%v", money*Rate, surMoney)
	//	if money*Rate <= surMoney {
	//		if p.ProgressBet >= 3 && p.ProgressBet <= 5 {
	//			betNum = 1
	//			data := &msg.ProgressBar_S2C{}
	//			data.ProBar = betNum
	//			p.SendMsg(data)
	//		} else if p.ProgressBet >= 6 && p.ProgressBet <= 8 {
	//			betNum = 2
	//			data := &msg.ProgressBar_S2C{}
	//			data.ProBar = betNum
	//			p.SendMsg(data)
	//		} else if p.ProgressBet >= 15 {
	//			betNum = 6
	//			// 发送进度条
	//			data := &msg.ProgressBar_S2C{}
	//			data.ProBar = betNum
	//			p.SendMsg(data)
	//			// 小游戏执行
	//			p.GetRewardsInfo()
	//		}
	//	} else { // 盈余池金额不足够小游戏获奖
	//		if p.ProgressBet >= 3 && p.ProgressBet <= 5 {
	//			betNum = 1
	//			data := &msg.ProgressBar_S2C{}
	//			data.ProBar = betNum
	//			p.SendMsg(data)
	//		} else if p.ProgressBet >= 6 {
	//			betNum = 2
	//			data := &msg.ProgressBar_S2C{}
	//			data.ProBar = betNum
	//			p.SendMsg(data)
	//		}
	//	}
	//}
}

func (p *Player) GodPickUpGold(betNum int32) {
	rid, _ := hall.UserRoom.Load(p.Id)
	v, _ := hall.RoomRecord.Load(rid)
	if v != nil {
		room := v.(*Room)

		// 获取财神接金币金额
		rate, money := GetGOLD(betNum)
		winMoney := money * CfgMoney[room.Config]

		// 结算
		pac := packageTax[p.PackageId]
		taxR := pac / 100
		tax := winMoney * taxR
		resultMoney := winMoney - tax

		p.Account += resultMoney
		p.TotalWinMoney += winMoney

		log.Debug("获取赢钱的金额:%v", winMoney)

		data := &msg.PickUpGold_S2C{}
		data.Money = winMoney
		data.Rate = rate
		data.Account = p.Account
		p.SendMsg(data)
	}
}

func (p *Player) HandleLuckyBag() {
	rid, _ := hall.UserRoom.Load(p.Id)
	v, _ := hall.RoomRecord.Load(rid)
	if v != nil {
		room := v.(*Room)

		rate := GetLuckyBag()
		surMoney := GetSurPlusMoney()
		luckyMoney := CfgMoney[room.Config] * float64(rate)
		if luckyMoney <= surMoney { // 判断福袋盈余池是否足够
			pac := packageTax[p.PackageId]
			taxR := pac / 100
			tax := luckyMoney * taxR
			resultMoney := luckyMoney - tax

			p.Account += resultMoney
			p.TotalWinMoney += luckyMoney

			data := &msg.LuckyBagAction_S2C{}
			data.IsDown = true
			data.Money = luckyMoney
			data.Account = p.Account
			p.SendMsg(data)
		}
	}
}

func (p *Player) ChangeRoomCfg(m *msg.ChangeRoomCfg_C2S) {
	// 判断玩家信息是否为空
	if p.Id == "" {
		p.SendErrMsg(RECODE_PlayerInfoIDIsNull)
		return
	}

	rid, _ := hall.UserRoom.Load(p.Id)
	v, _ := hall.RoomRecord.Load(rid)
	if v != nil {
		room := v.(*Room)
		// 保存区间节点位置
		p.ConfigPlace[m.Config] = m.Coordinates
		// 修改当前配置区间
		room.Config = m.ChangeCfg

		// 判断该金币区间是否存在金币位置存储，如果存在则返回，不存在则返回空
		if p.ConfigPlace[m.ChangeCfg] != nil {
			// 发送配置数据
			data := &msg.ChangeRoomCfg_S2C{}
			data.IsChange = true
			data.Coordinates = p.ConfigPlace[room.Config]
			p.SendMsg(data)
		} else {
			// 发送配置数据
			data := &msg.ChangeRoomCfg_S2C{}
			data.IsChange = false
			data.Coordinates = p.ConfigPlace[room.Config]
			p.SendMsg(data)
		}
	}
}
