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
	if v == nil {
		// 创建房间
		r := &Room{}
		r.Init()
		hall.RoomRecord.Store(r.RoomId, r)
		hall.UserRoom.Store(p.Id, r.RoomId)

		// 插入玩家信息 todo
		p.FindPlayerInfo()

		//返回前端房间信息
		data := &msg.JoinRoom_S2C{}
		data.RoomData = r.RespRoomData()
		p.SendMsg(data)
	}
}

func (p *Player) PlayerAction(downBet float64) {
	// 判断玩家金额是否足够
	if p.Account < downBet {
		log.Debug("玩家金额不足,不能进行下注~")
		return
	}
	// 判断下注金币是房间下注对应金额

	p.Account -= downBet

	p.LoseResultMoney = downBet
	nowTime := time.Now().Unix()
	p.RoundId = fmt.Sprintf("%+v-%+v", time.Now().Unix(), p.Id)
	loseReason := "发财推币机输钱"
	c2c.UserSyncLoseScore(p, nowTime, p.RoundId, loseReason, downBet)

	// 返回玩家行动数据
	action := &msg.PlayerAction_S2C{}
	action.DownBet = downBet
	action.Account = p.Account
	p.SendMsg(action)
}

func (p *Player) GetPlayerWinMoney(money float64) {
	p.Account += money

	p.WinResultMoney = money
	nowTime := time.Now().Unix()
	p.RoundId = fmt.Sprintf("%+v-%+v", time.Now().Unix(), p.Id)
	winReason := "发财推币机赢钱"
	c2c.UserSyncWinScore(p, nowTime, p.RoundId, winReason, money)

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

	// 判断是否大于跑马灯
}
