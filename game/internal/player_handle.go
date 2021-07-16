package internal

import (
	"fctbj/msg"
	"fmt"
	"github.com/name5566/leaf/log"
	"time"
)

func (p *Player) PlayerJoinRoom() {
	rid, _ := hall.UserRoom.Load(p.Id)
	v, _ := hall.RoomRecord.Load(rid)
	if v != nil { // 当前玩家已经存在房间
		room := v.(*Room)
		roomData := room.RespRoomData()
		enter := &msg.EnterRoom_S2C{}
		enter.RoomData = roomData
		p.SendMsg(enter)
	} else {
		// 创建房间
		r := &Room{}
		r.Init()
		hall.RoomRecord.Store(r.RoomId, r)
		hall.UserRoom.Store(p.Id, r.RoomId)

		// 插入玩家信息 todo
		p.FindPlayerInfo()

		log.Debug("发送进入房间!")

		//返回前端房间信息
		data := &msg.JoinRoom_S2C{}
		data.RoomData = r.RespRoomData()
		p.SendMsg(data)
	}
}

func (p *Player) PlayerAction(downBet float64) {
	// 判断玩家金额是否足够
	if p.Account < downBet {
		p.SendErrMsg(RECODE_UserMoneyNotEnough)
		log.Debug("玩家金额不足,不能进行下注~")
		return
	}

	// 判断下注金币是否对应房间配置金额(防止刷钱)
	var roomId string
	rid, _ := hall.UserRoom.Load(p.Id)
	v, _ := hall.RoomRecord.Load(rid)
	if v != nil {
		room := v.(*Room)
		roomId = room.RoomId
		if CfgMoney[room.Config] != downBet {
			log.Debug("房间配置金额:%v,%v", CfgMoney[room.Config], downBet)
			p.SendErrMsg(RECODE_RoomCfgMoneyERROR)
			log.Debug("房间配置金额不对!")
			return
		}
	}

	p.Account -= downBet
	p.LoseResultMoney = downBet

	nowTime := time.Now().Unix() // todo
	p.RoundId = fmt.Sprintf("%+v-%+v", time.Now().Unix(), p.Id)
	loseReason := "发财推币机输钱"
	c2c.UserSyncLoseScore(p, nowTime, p.RoundId, loseReason, downBet)

	// 插入盈余数据 todo
	sur := &SurplusPoolDB{}
	sur.UpdateTime = time.Now()
	sur.TimeNow = time.Now().Format("2006-01-02 15:04:05")
	sur.Rid = roomId
	sur.PlayerNum = LoadPlayerCount()
	surPool := FindSurplusPool()
	if surPool != nil {
		sur.HistoryWin = surPool.HistoryWin
		sur.HistoryLose = surPool.HistoryLose
	}
	sur.HistoryLose += Decimal(p.LoseResultMoney)
	sur.TotalLoseMoney += Decimal(p.LoseResultMoney)
	InsertSurplusPool(sur)

	// 先判断盈余池是否有钱，然后在处理玩家是否赢钱
	surMoney := GetSurPlusMoney()
	var goldNum int32
	if surMoney > 0 {
		goldNum = p.randGoldNum()
		if goldNum > 0 {
			// 玩家赢钱结算
			winMoney := downBet * float64(goldNum)
			pac := packageTax[p.PackageId]
			taxR := pac / 100
			tax := winMoney * taxR
			resultMoney := winMoney - tax

			p.Account += resultMoney
			p.WinResultMoney = winMoney

			log.Debug("获取赢钱的金额:%v", winMoney)

			nowTime := time.Now().Unix() // todo
			p.RoundId = fmt.Sprintf("%+v-%+v", time.Now().Unix(), p.Id)
			winReason := "发财推币机赢钱"
			c2c.UserSyncWinScore(p, nowTime, p.RoundId, winReason, winMoney)

			// 跑马灯 todo
			if resultMoney > PaoMaDeng {
				c2c.NoticeWinMoreThan(p.Id, p.NickName, resultMoney)
			}

			// 插入盈余数据 todo
			sur := &SurplusPoolDB{}
			sur.UpdateTime = time.Now()
			sur.TimeNow = time.Now().Format("2006-01-02 15:04:05")
			sur.Rid = roomId
			sur.PlayerNum = LoadPlayerCount()
			surPool := FindSurplusPool()
			if surPool != nil {
				sur.HistoryWin = surPool.HistoryWin
				sur.HistoryLose = surPool.HistoryLose
			}
			sur.HistoryWin += Decimal(p.WinResultMoney)
			sur.TotalWinMoney += Decimal(p.WinResultMoney)
			InsertSurplusPool(sur)
		}
	}

	data := &msg.PlayerAction_S2C{}
	data.WinNum = goldNum
	p.SendMsg(data)
}

func (p *Player) randGoldNum() int32 {
	var goldNum int32
	num := RandInRange(0, 100)
	if num >= 0 && num <= 65 {
		goldNum = 0
	} else if num >= 66 && num <= 85 {
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
		num := RandInRange(0, 100)
		if num >= 0 && num <= 5 {
			data.RewardsNum = GOLD
		} else if num >= 6 && num <= 12 {
			data.RewardsNum = RICH
			data.GetMoney = GetRICH(cfgMoney)
		} else if num >= 13 && num <= 30 {
			data.RewardsNum = PUSH
			data.GetMoney = GetPUSH(cfgMoney)
		} else if num >= 31 && num <= 100 {
			data.RewardsNum = LUCKY
			data.LuckyPig = GetLUCKY(cfgMoney)
		}
		p.SendMsg(data)
	}
}

func (p *Player) ProgressBetResp(bet int32) {
	p.ProgressBet += bet

	rid, _ := hall.UserRoom.Load(p.Id)
	v, _ := hall.RoomRecord.Load(rid)
	if v != nil {
		//room := v.(*Room)
		// 房间配置金额
		//money := CfgMoney[room.Config]
		//surMoney := GetSurPlusMoney()
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
		}

		//// 盈余池金额足够小游戏获奖时
		//if money*Rate >= surMoney {
		//	if p.ProgressBet >= 3 && p.ProgressBet <= 5 {
		//		betNum = 1
		//		data := &msg.ProgressBar_S2C{}
		//		data.ProBar = betNum
		//		p.SendMsg(data)
		//	} else if p.ProgressBet >= 6 && p.ProgressBet <= 8 {
		//		betNum = 2
		//		data := &msg.ProgressBar_S2C{}
		//		data.ProBar = betNum
		//		p.SendMsg(data)
		//	} else if p.ProgressBet >= 15 {
		//		betNum = 6
		//		// 发送进度条
		//		data := &msg.ProgressBar_S2C{}
		//		data.ProBar = betNum
		//		p.SendMsg(data)
		//		// 小游戏执行
		//		p.GetRewardsInfo()
		//	}
		//} else { // 盈余池金额不足够小游戏获奖
		//	if p.ProgressBet >= 3 && p.ProgressBet <= 5 {
		//		betNum = 1
		//		data := &msg.ProgressBar_S2C{}
		//		data.ProBar = betNum
		//		p.SendMsg(data)
		//	} else if p.ProgressBet >= 6 {
		//		betNum = 2
		//		data := &msg.ProgressBar_S2C{}
		//		data.ProBar = betNum
		//		p.SendMsg(data)
		//	}
		//}
	}
}

func (p *Player) ChangeRoomCfg(m *msg.ChangeRoomCfg_C2S) {
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
