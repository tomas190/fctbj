package internal

import (
	"fctbj/msg"
	"github.com/name5566/leaf/log"
	"strconv"
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
		log.Debug("玩家金额不足,不能进行下注:%v,%v", p.Account, downBet)
		return
	}

	var IsDown bool
	var coinName string
	rid, _ := hall.UserRoom.Load(p.Id)
	v, _ := hall.RoomRecord.Load(rid)
	if v != nil {
		room := v.(*Room)

		// 判断下注金币是否对应房间配置金额(防止刷钱)
		if CfgMoney[room.Config] != downBet {
			p.SendErrMsg(RECODE_RoomCfgMoneyERROR)
			log.Debug("房间配置金额错误:%v,%v", CfgMoney[room.Config], downBet)
			return
		}

		// 记录当前 Coin的序号 和 Coin列表
		room.CoinNum[room.Config]++
		coinName = Coin + strconv.Itoa(int(room.CoinNum[room.Config]))
		room.CoinList[room.Config] = append(room.CoinList[room.Config], coinName)

		// 判断是否掉落福袋
		p.DownBetCount++
		IsLucky := room.ExistLuckyBag()
		if p.DownBetCount >= 50 && IsLucky == false {
			IsDown = true
			p.DownBetCount = 0
			room.CoinList[room.Config] = append(room.CoinList[room.Config], FuDai)
		}
	}

	p.Account -= downBet
	p.DownBet += downBet
	p.TotalLoseMoney += downBet

	// 游戏赢率结算
	p.GameSurSettle()

	data := &msg.PlayerAction_S2C{}
	data.LuckyBag = IsDown
	data.Coin = coinName
	data.CoinList = p.DownBetList
	p.SendMsg(data)
}

func (p *Player) PlayerResult(coinList []string) {
	if coinList == nil {
		p.SendErrMsg(RECODE_ActionCoinNotHave)
		log.Debug("玩家行动金币为空!")
		return
	}

	rid, _ := hall.UserRoom.Load(p.Id)
	v, _ := hall.RoomRecord.Load(rid)
	if v != nil {
		room := v.(*Room)
		// 获取相同的金币进行赢钱结算
		var winNum int
		var luckyBag bool
		for _, v := range coinList {
			if v == FuDai {
				luckyBag = true
			}
			// 判断获取相同的金币并删除
			for k, c := range room.CoinList[room.Config] {
				if v == c {
					winNum++
					room.CoinList[room.Config] = append(room.CoinList[room.Config][:k], room.CoinList[room.Config][k+1:]...)
				}
			}
		}

		// 玩家赢钱结算
		var winMoney float64
		// 福袋结算
		if luckyBag == true {
			winMoney += CfgMoney[room.Config] * float64(LuckyBag)
		}
		// 金币结算
		winMoney += CfgMoney[room.Config] * float64(winNum)
		pac := packageTax[p.PackageId]
		taxR := pac / 100
		tax := winMoney * taxR
		resultMoney := winMoney - tax

		p.Account += resultMoney
		p.TotalWinMoney += winMoney
		log.Debug("获取赢钱的金额:%v", winMoney)

		data := &msg.ActionResult_S2C{}
		data.Account = p.Account
		p.SendMsg(data)
	}
}

func (p *Player) GameSurSettle() {
	sur := GetFindSurPool()
	loseRate := sur.PlayerLoseRateAfterSurplusPool * 100
	percentageWin := sur.RandomPercentageAfterWin * 100
	percentageLose := sur.RandomPercentageAfterLose * 100
	countWin := sur.RandomCountAfterWin
	countLose := sur.RandomCountAfterLose
	surplusPool := sur.SurplusPool

	num := RandInRange(1, 101)
	if num >= 0 { // 玩家赢钱
		settle := p.GetGoldSettle()
		for {
			loseRateNum := RandInRange(1, 101)
			percentageWinNum := RandInRange(1, 101)
			if countWin > 0 {
				if percentageWinNum > int(percentageWin) { // 盈余池判定
					if surplusPool > settle { // 盈余池足够
						break
					} else {                             // 盈余池不足
						if loseRateNum > int(loseRate) { // 30%玩家赢钱
							break
						} else { // 70%玩家输钱
							p.DownBetList = nil
							break
						}
					}
				} else { // 又随机生成牌型
					settle := p.GetGoldSettle()
					if settle > 0 { // 玩家赢
						countWin--
					} else {
						break
					}
				}
			} else {
				// 盈余池判定
				if surplusPool > settle { // 盈余池足够
					break
				} else {                             // 盈余池不足
					if loseRateNum > int(loseRate) { // 30%玩家赢钱
						break
					} else { // 70%玩家输钱
						p.DownBetList = nil
						return
					}
				}
			}
		}
	} else { // 玩家输钱
		for {
			loseRateNum := RandInRange(1, 101)
			percentageLoseNum := RandInRange(1, 101)
			if countLose > 0 {
				if percentageLoseNum > int(percentageLose) {
					break
				} else { // 又随机生成牌型
					settle := p.GetGoldSettle()
					if settle > 0 { // 玩家赢
						// 盈余池判定
						if surplusPool > settle { // 盈余池足够
							break
						} else {                             // 盈余池不足
							if loseRateNum > int(loseRate) { // 30%玩家赢钱
								for {
									settle := p.GetGoldSettle()
									if settle >= 0 {
										return
									}
								}
							} else { // 70%玩家输钱
								p.DownBetList = nil
								return
							}
						}
					} else {
						countLose--
					}
				}
			} else { // 玩家输钱
				p.DownBetList = nil
				return
			}
		}
	}
}

func (p *Player) randGoldNum() int {
	num := RandInRange(1, 101)
	return num
}

func (p *Player) GetGoldSettle() float64 {
	rid, _ := hall.UserRoom.Load(p.Id)
	v, _ := hall.RoomRecord.Load(rid)
	if v != nil {
		room := v.(*Room)
		if len(room.CoinList[room.Config]) == 0 {
			p.SendErrMsg(RECODE_TableNotHaveGold)
			log.Debug("玩家桌面金币为空!")
			return 0
		}
		// 房间配置金额
		cfgMoney := CfgMoney[room.Config]
		var goldNum int
		for { // 循环获取随机金币,避免随机到金币大于桌面金币数量
			goldNum = p.randGoldNum()
			if len(room.CoinList[room.Config]) >= goldNum {
				break
			}
		}
		p.DownBetList = room.CoinList[room.Config][:goldNum]
		log.Debug("随机金币的数量:%v,切片值:%v", goldNum, p.DownBetList)
		settle := cfgMoney * float64(goldNum)
		return settle
	}
	return 0
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
		var fudai1 int
		var fudai2 int
		num := RandInRange(0, 100)
		if num >= 0 && num <= 5 {
			data.RewardsNum = GOLD
		} else if num >= 6 && num <= 12 {
			data.RewardsNum = RICH
			data.GetMoney = GetRICH(cfgMoney)
			winMoney = data.GetMoney
		} else if num >= 13 && num <= 30 {
			data.RewardsNum = PUSH
			winMoney, fudai1, fudai2 = room.GetPUSH(cfgMoney)
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

		down := &msg.DownLuckyBag_S2C{}
		down.LuckyBag1 = int32(fudai1)
		down.LuckyBag2 = int32(fudai2)
		down.CoinList = room.CoinList[room.Config]
		down.Money = resultMoney
		p.SendMsg(down)

		// Push中奖,清除桌面金币和福袋,重新生成新的金币
		if data.RewardsNum == PUSH {
			room.CoinList[room.Config] = nil
			for i := 1; i <= 100; i++ {
				room.CoinNum[room.Config] ++
				room.CoinList[room.Config] = append(room.CoinList[room.Config], Coin+strconv.Itoa(int(room.CoinNum[room.Config])))
			}
			creat := &msg.ReCreatGold_S2C{}
			creat.CoinList = room.CoinList[room.Config]
			p.SendMsg(creat)
		}
	}
}

func (p *Player) ProgressBetResp(bet int32, coin string) {

	var betNum int32
	rid, _ := hall.UserRoom.Load(p.Id)
	v, _ := hall.RoomRecord.Load(rid)
	if v != nil {
		room := v.(*Room)

		p.ProgressBet += bet
		log.Debug("p.ProgressBet 长度为:%v", p.ProgressBet)

		for k, v := range room.CoinList[room.Config] {
			if v == coin {
				room.CoinList[room.Config] = append(room.CoinList[room.Config][:k], room.CoinList[room.Config][k+1:]...)
			}
		}

		// 房间配置金额
		money := CfgMoney[room.Config]
		surMoney := GetSurPlusMoney()
		// 盈余池金额足够小游戏获奖时
		log.Debug("获奖的估计金额:%v,盈余池金额:%v", money*Rate, surMoney)
		if money*Rate <= surMoney {
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
			} else if p.ProgressBet >= 15 {
				betNum = 6
				// 发送进度条
				data := &msg.ProgressBar_S2C{}
				data.ProBar = betNum
				p.SendMsg(data)
				// 小游戏执行
				p.GetRewardsInfo()
			}
		} else { // 盈余池金额不足够小游戏获奖
			if p.ProgressBet >= 3 && p.ProgressBet <= 5 {
				betNum = 1
				data := &msg.ProgressBar_S2C{}
				data.ProBar = betNum
				p.SendMsg(data)
			} else if p.ProgressBet >= 6 {
				betNum = 2
				data := &msg.ProgressBar_S2C{}
				data.ProBar = betNum
				p.SendMsg(data)
			}
		}
	}
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
		data.Money = resultMoney
		data.Rate = rate
		data.Account = p.Account
		p.SendMsg(data)
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
			data.CoinList = room.CoinList[room.Config]
			data.Coordinates = p.ConfigPlace[room.Config]
			p.SendMsg(data)
		} else {
			// 发送配置数据
			data := &msg.ChangeRoomCfg_S2C{}
			data.IsChange = false
			data.CoinList = room.CoinList[room.Config]
			data.Coordinates = p.ConfigPlace[room.Config]
			p.SendMsg(data)
		}
	}
}
