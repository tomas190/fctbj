package internal

import (
	"fctbj/msg"
	"fmt"
	"github.com/name5566/leaf/log"
	"strconv"
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

func (p *Player) PlayerAction(bet string) {
	downBet, _ := strconv.ParseFloat(bet, 64)
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

	var isWin bool
	surMoney := GetSurPlusMoney()
	if surMoney > downBet {
		isWin = true
	}

	data := &msg.PlayerAction_S2C{}
	data.IsWin = isWin
	p.SendMsg(data)
}

func (p *Player) GetPlayerWinMoney(bet string) {
	money, _ := strconv.ParseFloat(bet, 64)
	if money <= 0 {
		p.SendErrMsg(RECODE_SendWinMoneyERROR)
		log.Debug("玩家赢钱金额小于0 错误!")
		return
	}

	pac := packageTax[p.PackageId]
	taxR := pac / 100
	tax := money * taxR
	resultMoney := money - tax
	p.Account += resultMoney

	log.Debug("获取赢钱的金额:%v", money)
	p.WinResultMoney = money

	nowTime := time.Now().Unix() // todo
	p.RoundId = fmt.Sprintf("%+v-%+v", time.Now().Unix(), p.Id)
	winReason := "发财推币机赢钱"
	c2c.UserSyncWinScore(p, nowTime, p.RoundId, winReason, money)

	// 插入盈余数据 todo
	sur := &SurplusPoolDB{}
	sur.UpdateTime = time.Now()
	sur.TimeNow = time.Now().Format("2006-01-02 15:04:05")
	rid, _ := hall.UserRoom.Load(p.Id)
	sur.Rid = rid.(string)
	sur.PlayerNum = LoadPlayerCount()
	surPool := FindSurplusPool()
	if surPool != nil {
		sur.HistoryWin = surPool.HistoryWin
		sur.HistoryLose = surPool.HistoryLose
	}
	sur.HistoryLose += Decimal(p.WinResultMoney)
	sur.TotalLoseMoney += Decimal(p.WinResultMoney)
	InsertSurplusPool(sur)

	// 跑马灯 todo
	if resultMoney > PaoMaDeng {
		c2c.NoticeWinMoreThan(p.Id, p.NickName, resultMoney)
	}

	data := &msg.SendWinMoney_S2S{}
	data.Account = p.Account
	p.SendMsg(data)
}

func (p *Player) GetRewardsInfo() {
	var rewardsNum int32
	var rewardsMoney int32
	num := RandInRange(0, 100)
	if num >= 0 && num <= 5 {
		rewardsNum = GOLD
	} else if num >= 6 && num <= 12 {
		rewardsNum = RICH
	} else if num >= 13 && num <= 30 {
		rewardsNum = PUSH
	} else if num >= 31 && num <= 100 {
		rewardsNum = LUCKY
	}

	num2 := RandInRange(0, 100)
	if num2 >= 0 && num2 <= 20 {
		rewardsMoney = 1
	} else if num2 >= 21 && num2 <= 50 {
		rewardsMoney = 2
	} else if num2 >= 51 && num2 <= 100 {
		rewardsMoney = 3
	}

	data := &msg.GetRewards_S2C{}
	data.RewardsNum = rewardsNum
	data.RewardsMoney = rewardsMoney
	p.SendMsg(data)

}
