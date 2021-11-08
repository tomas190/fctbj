package internal

import (
	"fctbj/conf"
	"fctbj/msg"
	"github.com/name5566/leaf/log"
	"strconv"
	"time"
)

func (p *Player) PlayerJoinRoom(cfgId string) {
	// 判断玩家信息是否为空
	if p.Id == "" {
		p.SendErrMsg(RECODE_PlayerInfoIDIsNull)
		return
	}
	if p.IsExist == true {
		p.SendErrMsg(RECODE_PlayerExistRoom)
		return
	}

	// 创建房间
	r := &Room{}
	r.Init()
	r.Config = cfgId
	hall.RoomRecord.Store(r.RoomId, r)
	hall.UserRoom.Store(p.Id, r.RoomId)

	p.RoomId = r.RoomId
	p.IsExist = true

	// 插入玩家信息 todo
	p.FindPlayerInfo()

	//返回前端房间信息
	data := &msg.JoinRoom_S2C{}
	data.RoomData = r.RespRoomData()
	p.SendMsg(data)
}

func (p *Player) ExitFromRoom(room *Room) {
	// 判断玩家期间是否行动过
	//if p.DownBet > 0 {
	//	// 插入玩家数据
	//	p.HandlePlayerData()
	//}

	p.IsExist = false
	p.OffLineTime = time.Now().Hour()
	p.ProgressBet = 0
	p.DownBetCount = 0
	hall.OnlineUser.Delete(p.Id)

	c2c.UserLogoutCenter(p.Id, p.Password, p.Token) //todo
	leaveHall := &msg.Logout_S2C{}
	p.SendMsg(leaveHall)
	p.ConnAgent.Close()
}

func (p *Player) PlayerAction(m *msg.PlayerAction_C2S) {
	// 判断玩家信息是否为空
	if p.Id == "" {
		p.SendErrMsg(RECODE_PlayerInfoIDIsNull)
		log.Debug("玩家信息ID为空,不能进行下注")
		return
	}
	// 判断玩家金额是否足够
	if p.Account < m.DownBet {
		p.SendErrMsg(RECODE_UserMoneyNotEnough)
		log.Debug("玩家金额不足,不能进行下注:%v,%v", p.Account, m.DownBet)
		return
	}
	// 检索玩家是否为登出状态
	if p.IsLogin == false {
		log.Debug("玩家为登出状态,行动失败!")
		return
	}
	if p.IsExist == false {
		log.Debug("当前玩家当前已退出状态,行动失败!")
		return
	}

	var IsDown bool
	var coinName string
	var storageCoin []string
	rid, _ := hall.UserRoom.Load(p.Id)
	v, _ := hall.RoomRecord.Load(rid)
	if v != nil {
		room := v.(*Room)
		if room.IsLuckyGame == true {
			p.SendErrMsg(RECODE_InRoomGameStep)
			log.Debug("玩家行动失败,房间正在小游戏!")
			return
		}
		// 判断下注金币是否对应房间配置金额(防止刷钱)
		if CfgMoney[room.Config] != m.DownBet {
			p.SendErrMsg(RECODE_RoomCfgMoneyERROR)
			log.Debug("房间配置金额错误:%v,%v", CfgMoney[room.Config], m.DownBet)
			return
		}

		log.Debug("当前房间配置:%v,玩家下注金额:%v", room.Config, m.DownBet)

		// 记录当前 Coin的序号 和 Coin列表
		room.CoinNum[room.Config]++
		coinName = Coin + strconv.Itoa(int(room.CoinNum[room.Config]))
		room.CoinList[room.Config] = append(room.CoinList[room.Config], coinName)

		//log.Debug("金币长度:%v,金币列表:%v", len(room.CoinList[room.Config]), room.CoinList[room.Config])
		// 判断是否掉落福袋
		p.DownBetCount++
		IsLucky := room.ExistLuckyBag()
		if p.DownBetCount >= 50 && IsLucky == false {
			IsDown = true
			p.DownBetCount = 0
			room.CoinList[room.Config] = append(room.CoinList[room.Config], FuDai)
			log.Debug("福袋掉落:%v,%v,长度:%v", room.Config, room.CoinList[room.Config], len(room.CoinList[room.Config]))
		}
		storageCoin = room.CoinList[room.Config]
	}
	p.DownBet = 0

	p.Account -= m.DownBet
	p.LoseResultMoney = m.DownBet
	p.DownBet += m.DownBet
	p.TotalLoseMoney += m.DownBet

	// todo
	nowTime := time.Now().Unix()
	p.RoundId = p.RandRoundId()
	loseReason := "发财推币机输钱"
	c2c.UserSyncLoseScore(p, nowTime, p.RoundId, loseReason, m.DownBet)

	// 游戏赢率结算
	p.GameSurSettle()

	log.Debug("玩家行动获奖长度:%v", len(p.DownBetList))

	data := &msg.PlayerAction_S2C{}
	data.LuckyBag = IsDown
	data.Coin = coinName
	data.CoinList = p.DownBetList
	data.StorageList = storageCoin
	p.SendMsg(data)

	pac := packageTax[p.PackageId]
	taxR := pac / 100

	go func() {
		// 插入运营数据
		pr := &PlayerDownBetRecode{}
		pr.Id = p.Id
		pr.GameId = conf.Server.GameID
		pr.RoundId = p.RoundId
		pr.RoomId = p.RoomId
		pr.DownBetInfo = p.DownBet
		pr.DownBetTime = nowTime
		pr.StartTime = nowTime
		pr.EndTime = nowTime
		pr.SettlementFunds = p.LoseResultMoney
		pr.SpareCash = p.Account
		pr.TaxRate = taxR
		InsertAccessData(pr)

		// 插入游戏统计数据
		sd := &StatementData{}
		sd.Id = p.Id
		sd.GameId = conf.Server.GameID
		sd.GameName = "财神推金币"
		sd.DownBetTime = nowTime
		sd.StartTime = nowTime
		sd.EndTime = nowTime
		sd.PackageId = p.PackageId
		sd.LoseStatementTotal = p.LoseResultMoney
		sd.BetMoney = p.DownBet
		InsertStatementDB(sd)

		if p.PackageId != 11 && p.PackageId != 8 {
			// 插入盈余数据
			sur := &SurplusPoolDB{}
			sur.UpdateTime = time.Now()
			sur.TimeNow = time.Now().Format("2006-01-02 15:04:05")
			sur.Rid = p.RoomId
			sur.PlayerNum = LoadPlayerCount()
			surPool := FindSurplusPool()
			if surPool != nil {
				sur.HistoryWin = surPool.HistoryWin
				sur.HistoryLose = surPool.HistoryLose
			}
			sur.HistoryLose += Decimal(p.LoseResultMoney)
			sur.TotalLoseMoney += Decimal(p.LoseResultMoney)
			InsertSurplusPool(sur)
		}
	}()
}

func (p *Player) PlayerResult(m *msg.ActionResult_C2S) {
	if m.CoinList == nil {
		p.SendErrMsg(RECODE_ActionCoinNotHave)
		log.Debug("玩家结算金币为空!")
		return
	}

	rid, _ := hall.UserRoom.Load(p.Id)
	v, _ := hall.RoomRecord.Load(rid)
	if v != nil {
		room := v.(*Room)
		//if room.IsLuckyGame == true {
		//	p.SendErrMsg(RECODE_InRoomGameStep)
		//	log.Debug("玩家结算失败,房间正在小游戏!")
		//	return
		//}
		// 获取相同的金币进行赢钱结算
		var winNum int
		var luckyBag bool
		for _, v := range m.CoinList {
			if v == FuDai {
				luckyBag = true
				log.Debug("福袋结算:%v", v)
			}
			// 判断获取相同的金币并删除
			for k, c := range room.CoinList[room.Config] {
				if v == c {
					winNum++
					room.CoinList[room.Config] = append(room.CoinList[room.Config][:k], room.CoinList[room.Config][k+1:]...)
				}
			}
		}

		if winNum == 0 && luckyBag == false {
			p.SendErrMsg(RECODE_NotHaveSameCoin)
			log.Debug("金币列表未存在相同的金币:%v,%v", m.CoinList, room.CoinList[room.Config])
			return
		}

		log.Debug("===>结算当前金币数量长度为:%v", len(room.CoinList[room.Config]))

		// 玩家赢钱结算
		var winMoney float64
		// 福袋结算
		if luckyBag == true {
			winNum--
			winMoney += CfgMoney[room.Config] * float64(LuckyBag)
			log.Debug("福袋结算列表:%v", room.CoinList[room.Config])
		}
		// 金币结算
		winMoney += CfgMoney[room.Config] * float64(winNum)
		pac := packageTax[p.PackageId]
		taxR := pac / 100
		tax := winMoney * taxR
		resultMoney := winMoney - tax

		p.Account += resultMoney
		p.WinResultMoney = winMoney
		p.TotalWinMoney += winMoney
		log.Debug("获取赢钱的金额:%v", winMoney)

		// todo
		nowTime := time.Now().Unix()
		p.RoundId = p.RandRoundId()
		winReason := "发财推币机赢钱"
		c2c.UserSyncWinScore(p, nowTime, p.RoundId, winReason, winMoney)

		// 跑马灯
		if resultMoney > PaoMaDeng {
			c2c.NoticeWinMoreThan(p.Id, p.NickName, resultMoney)
		}

		data := &msg.ActionResult_S2C{}
		data.Account = p.Account
		p.SendMsg(data)

		go func() {
			// 插入运营数据
			pr := &PlayerDownBetRecode{}
			pr.Id = p.Id
			pr.GameId = conf.Server.GameID
			pr.RoundId = p.RoundId
			pr.RoomId = p.RoomId
			pr.DownBetInfo = p.DownBet
			pr.DownBetTime = nowTime
			pr.StartTime = nowTime
			pr.EndTime = nowTime
			pr.SettlementFunds = resultMoney
			pr.SpareCash = p.Account
			pr.TaxRate = taxR
			InsertAccessData(pr)

			// 插入游戏统计数据
			sd := &StatementData{}
			sd.Id = p.Id
			sd.GameId = conf.Server.GameID
			sd.GameName = "财神推金币"
			sd.DownBetTime = nowTime
			sd.StartTime = nowTime
			sd.EndTime = nowTime
			sd.PackageId = p.PackageId
			sd.WinStatementTotal = p.WinResultMoney
			sd.BetMoney = p.DownBet
			InsertStatementDB(sd)

			if p.PackageId != 11 && p.PackageId != 8 {
				// 插入盈余数据
				sur := &SurplusPoolDB{}
				sur.UpdateTime = time.Now()
				sur.TimeNow = time.Now().Format("2006-01-02 15:04:05")
				sur.Rid = p.RoomId
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
		}()
	}
}

func (p *Player) GameSurSettle() {
	sur := GetFindSurPool() //todo
	loseRate := sur.PlayerLoseRateAfterSurplusPool * 100
	percentageWin := sur.RandomPercentageAfterWin * 100
	percentageLose := sur.RandomPercentageAfterLose * 100
	countWin := sur.RandomCountAfterWin
	countLose := sur.RandomCountAfterLose
	surplusPool := sur.SurplusPool

	if p.PackageId == 11 || p.PackageId == 8 {
		percentageWin = 38
		percentageLose = 0
		countWin = 3
		countLose = 0
	}

	num := RandInRange(1, 101)
	if num >= 50 { // 玩家赢钱
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
			num := RandInRange(1, 101)
			if num >= 55 {
				goldNum = p.randGoldNum()
				if len(room.CoinList[room.Config]) >= goldNum {
					break
				}
			} else {
				p.DownBetList = nil
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
		p.ProgressBet = 0
		room.IsLuckyGame = true

		// 房间配置金额
		cfgMoney := CfgMoney[room.Config]

		var winMoney float64
		var fudai1 int
		var fudai2 int
		var gameName string
		var rate float64

		data := &msg.GetRewards_S2C{}
		if len(room.CoinList[room.Config]) > 200 {
			data.RewardsNum = PUSH
			gameName = "财神发钱"
			rate, winMoney, fudai1, fudai2 = room.GetPUSH(cfgMoney)
		} else {
			num := RandInRange(1, 101)
			if num >= 1 && num <= 3 {
				data.RewardsNum = GOLD
				room.IsPickGod = true
			} else if num >= 4 && num <= 6 {
				data.RewardsNum = RICH
				gameName = "金猪送财"
				rate, winMoney = GetRICH(cfgMoney)
			} else if num >= 7 && num <= 10 {
				data.RewardsNum = PUSH
				gameName = "财神发钱"
				rate, winMoney, fudai1, fudai2 = room.GetPUSH(cfgMoney)
			} else if num >= 11 && num <= 100 {
				data.RewardsNum = LUCKY
				room.IsLuckyPig = true
			}
		}
		p.SendMsg(data)

		// 结算
		pac := packageTax[p.PackageId]
		taxR := pac / 100
		tax := winMoney * taxR
		resultMoney := winMoney - tax
		log.Debug("获取小游戏赢钱的金额:%v,%v,%v", gameName, winMoney, rate)

		p.Account += resultMoney
		p.WinResultMoney = winMoney
		p.TotalWinMoney += winMoney

		// todo
		nowTime := time.Now().Unix()
		if winMoney > 0 {
			winReason := "发财推币机" + gameName + "赢钱"
			p.RoundId = p.RandRoundId()
			c2c.UserSyncWinScore(p, nowTime, p.RoundId, winReason, winMoney)
		}

		// 跑马灯
		if resultMoney > PaoMaDeng {
			c2c.NoticeWinMoreThan(p.Id, p.NickName, resultMoney)
		}

		if data.RewardsNum == RICH {
			send := &msg.SendMoney_S2C{}
			send.GetMoney = resultMoney
			send.Account = p.Account
			p.SendMsg(send)

			go func() {
				timeout := time.NewTimer(time.Second * 5)
				for {
					select {
					case <-timeout.C:
						room.IsLuckyGame = false
						return
					}
				}
			}()
		}

		// Push中奖,清除桌面金币和福袋,重新生成新的金币
		if data.RewardsNum == PUSH {
			// 清空累计下注次数
			p.DownBetCount = 0

			down := &msg.DownLuckyBag_S2C{}
			down.LuckyBag1 = int32(fudai1)
			down.LuckyBag2 = int32(fudai2)
			down.CoinList = room.CoinList[room.Config]
			down.Money = resultMoney
			down.Account = p.Account
			p.SendMsg(down)

			room.CoinList[room.Config] = nil
			room.ConfigPlace[room.Config] = nil
			for i := 1; i <= 100; i++ {
				room.CoinNum[room.Config] ++
				room.CoinList[room.Config] = append(room.CoinList[room.Config], Coin+strconv.Itoa(int(room.CoinNum[room.Config])))
			}
			log.Debug("Push游戏生成金币长度:%v", len(room.CoinList[room.Config]))
			creat := &msg.ReCreatGold_S2C{}
			creat.CoinList = room.CoinList[room.Config]
			p.SendMsg(creat)

			go func() {
				timeout := time.NewTimer(time.Second * 5)
				for {
					select {
					case <-timeout.C:
						room.IsLuckyGame = false
						return
					}
				}
			}()
		}

		go func() {
			if p.WinResultMoney > 0 {
				// 插入运营数据
				pr := &PlayerDownBetRecode{}
				pr.Id = p.Id
				pr.GameId = conf.Server.GameID
				pr.RoundId = p.RoundId
				pr.RoomId = p.RoomId
				pr.DownBetInfo = p.DownBet
				pr.DownBetTime = nowTime
				pr.StartTime = nowTime
				pr.EndTime = nowTime
				pr.GameReward = new(GameRewards)
				pr.GameReward.Game = gameName
				pr.GameReward.Rate = rate
				pr.GameReward.WinMoney = winMoney
				pr.SettlementFunds = resultMoney
				pr.SpareCash = p.Account
				pr.TaxRate = taxR
				InsertAccessData(pr)

				// 插入游戏统计数据
				sd := &StatementData{}
				sd.Id = p.Id
				sd.GameId = conf.Server.GameID
				sd.GameName = "财神推金币"
				sd.DownBetTime = nowTime
				sd.StartTime = nowTime
				sd.EndTime = nowTime
				sd.PackageId = p.PackageId
				sd.WinStatementTotal = p.WinResultMoney
				sd.BetMoney = p.DownBet
				InsertStatementDB(sd)

				if p.PackageId != 11 && p.PackageId != 8 {
					// 插入盈余数据
					sur := &SurplusPoolDB{}
					sur.UpdateTime = time.Now()
					sur.TimeNow = time.Now().Format("2006-01-02 15:04:05")
					sur.Rid = p.RoomId
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
		}()
	}
}

func (p *Player) ProgressBetResp(m *msg.ProgressBar_C2S) {
	rid, _ := hall.UserRoom.Load(p.Id)
	v, _ := hall.RoomRecord.Load(rid)
	if v != nil {
		room := v.(*Room)
		if m.Coin == "" {
			p.SendErrMsg(RECODE_ProBarCoinNotHave)
			log.Debug("获取进度条金币为空:%v", m)
			return
		}

		p.ProgressBet++
		log.Debug("p.ProgressBet 进度条数量为:%v", p.ProgressBet)

		for k, v := range room.CoinList[room.Config] {
			if v == m.Coin {
				if v == FuDai {
					log.Debug("进度条福袋掉落~")
				}
				room.CoinList[room.Config] = append(room.CoinList[room.Config][:k], room.CoinList[room.Config][k+1:]...)
			}
		}
		log.Debug("===>进度条金币列表长度:%v", len(room.CoinList[room.Config]))

		// 房间配置金额
		money := CfgMoney[room.Config]
		surMoney := GetSurPlusMoney() //todo
		// 盈余池金额足够小游戏获奖时
		//log.Debug("获奖的估计金额:%v,盈余池金额:%v", money*Rate, surMoney)
		var betNum int32
		if money*Rate <= surMoney {
			if p.ProgressBet >= 1 && p.ProgressBet <= 50 {
				betNum = 1
				data := &msg.ProgressBar_S2C{}
				data.ProBar = betNum
				data.Coordinates = room.ConfigPlace[room.Config]
				p.SendMsg(data)
			} else if p.ProgressBet >= 51 && p.ProgressBet <= 99 {
				betNum = 2
				data := &msg.ProgressBar_S2C{}
				data.ProBar = betNum
				data.Coordinates = room.ConfigPlace[room.Config]
				p.SendMsg(data)
			} else if p.ProgressBet >= 100 {
				betNum = 6
				// 发送进度条
				data := &msg.ProgressBar_S2C{}
				data.ProBar = betNum
				data.Coordinates = room.ConfigPlace[room.Config]
				p.SendMsg(data)
				// 小游戏执行
				p.GetRewardsInfo()
			}
		} else { // 盈余池金额不足够小游戏获奖
			if p.ProgressBet >= 1 && p.ProgressBet <= 50 {
				betNum = 1
				data := &msg.ProgressBar_S2C{}
				data.ProBar = betNum
				data.Coordinates = room.ConfigPlace[room.Config]
				p.SendMsg(data)
			} else if p.ProgressBet >= 51 {
				betNum = 2
				data := &msg.ProgressBar_S2C{}
				data.ProBar = betNum
				data.Coordinates = room.ConfigPlace[room.Config]
				p.SendMsg(data)
			}
		}
	}
}

func (p *Player) WinLuckyPig() {
	rid, _ := hall.UserRoom.Load(p.Id)
	v, _ := hall.RoomRecord.Load(rid)
	if v != nil {
		room := v.(*Room)
		if room.IsActPig == true { // 防止多次点击请求
			return
		}
		if room.IsLuckyGame == false {
			log.Debug("当前不为小游戏行动阶段!")
			return
		}

		room.IsLuckyPig = false
		room.IsActPig = true
		var getPigInfo *msg.ThreePig
		// 获取财神接金币金额
		cfgMoney := CfgMoney[room.Config]
		rate, getPigInfo := GetLUCKY(cfgMoney)
		winMoney := getPigInfo.PigSuccess

		// 结算
		pac := packageTax[p.PackageId]
		taxR := pac / 100
		tax := winMoney * taxR
		resultMoney := winMoney - tax

		p.Account += resultMoney
		p.WinResultMoney = winMoney
		p.TotalWinMoney += winMoney

		log.Debug("财运满满赢钱的金额:%v,%v", winMoney, rate)

		// todo
		nowTime := time.Now().Unix()
		if winMoney > 0 {
			winReason := "发财推币机财运满满赢钱"
			p.RoundId = p.RandRoundId()
			c2c.UserSyncWinScore(p, nowTime, p.RoundId, winReason, winMoney)
		}

		// 跑马灯
		if resultMoney > PaoMaDeng {
			c2c.NoticeWinMoreThan(p.Id, p.NickName, resultMoney)
		}

		data := &msg.LuckyPig_S2C{}
		data.LuckyPig = getPigInfo
		data.Account = p.Account
		p.SendMsg(data)

		go func() {
			timeout := time.NewTimer(time.Second * 3)
			for {
				select {
				case <-timeout.C:
					room.IsActPig = false
					return
				}
			}
		}()

		room.IsLuckyGame = false

		go func() {
			if p.WinResultMoney > 0 { // todo
				// 插入运营数据
				pr := &PlayerDownBetRecode{}
				pr.Id = p.Id
				pr.GameId = conf.Server.GameID
				pr.RoundId = p.RoundId
				pr.RoomId = p.RoomId
				pr.DownBetInfo = p.DownBet
				pr.DownBetTime = nowTime
				pr.StartTime = nowTime
				pr.EndTime = nowTime
				pr.GameReward = new(GameRewards)
				pr.GameReward.Game = "财运满满"
				pr.GameReward.Rate = rate
				pr.GameReward.WinMoney = winMoney
				pr.SettlementFunds = resultMoney
				pr.SpareCash = p.Account
				pr.TaxRate = taxR
				InsertAccessData(pr)

				// 插入游戏统计数据
				sd := &StatementData{}
				sd.Id = p.Id
				sd.GameId = conf.Server.GameID
				sd.GameName = "财神推金币"
				sd.DownBetTime = nowTime
				sd.StartTime = nowTime
				sd.EndTime = nowTime
				sd.PackageId = p.PackageId
				sd.WinStatementTotal = p.WinResultMoney
				sd.BetMoney = p.DownBet
				InsertStatementDB(sd)

				if p.PackageId != 11 && p.PackageId != 8 {
					// 插入盈余数据
					sur := &SurplusPoolDB{}
					sur.UpdateTime = time.Now()
					sur.TimeNow = time.Now().Format("2006-01-02 15:04:05")
					sur.Rid = p.RoomId
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
			log.Debug("财运满满执行完毕!")
		}()
	}
}

func (p *Player) GodPickUpGold(betNum int32) {
	rid, _ := hall.UserRoom.Load(p.Id)
	v, _ := hall.RoomRecord.Load(rid)
	if v != nil {
		room := v.(*Room)
		if room.IsLuckyGame == false {
			log.Debug("当前不为小游戏行动阶段!")
			return
		}

		// 接金币小游戏结束
		room.IsPickGod = false

		// 获取财神接金币金额
		rate, money := GetGOLD(betNum)
		winMoney := money * CfgMoney[room.Config]

		// 结算
		pac := packageTax[p.PackageId]
		taxR := pac / 100
		tax := winMoney * taxR
		resultMoney := winMoney - tax

		p.Account += resultMoney
		p.WinResultMoney = winMoney
		p.TotalWinMoney += winMoney

		log.Debug("财神接金币赢钱的金额:%v", winMoney)

		// todo
		nowTime := time.Now().Unix()
		if winMoney > 0 {
			winReason := "发财推币机财神接金币赢钱"
			p.RoundId = p.RandRoundId()
			c2c.UserSyncWinScore(p, nowTime, p.RoundId, winReason, winMoney)
		}

		// 跑马灯
		if resultMoney > PaoMaDeng {
			c2c.NoticeWinMoreThan(p.Id, p.NickName, resultMoney)
		}

		data := &msg.PickUpGold_S2C{}
		data.Money = resultMoney
		data.Rate = rate
		data.Account = p.Account
		p.SendMsg(data)

		room.IsLuckyGame = false

		go func() {
			if p.WinResultMoney > 0 {
				// 插入运营数据
				pr := &PlayerDownBetRecode{}
				pr.Id = p.Id
				pr.GameId = conf.Server.GameID
				pr.RoundId = p.RoundId
				pr.RoomId = p.RoomId
				pr.DownBetInfo = p.DownBet
				pr.DownBetTime = nowTime
				pr.StartTime = nowTime
				pr.EndTime = nowTime
				pr.GameReward = new(GameRewards)
				pr.GameReward.Game = "财神接金币"
				pr.GameReward.Rate = float64(rate)
				pr.GameReward.WinMoney = winMoney
				pr.SettlementFunds = resultMoney
				pr.SpareCash = p.Account
				pr.TaxRate = taxR
				InsertAccessData(pr)

				// 插入游戏统计数据
				sd := &StatementData{}
				sd.Id = p.Id
				sd.GameId = conf.Server.GameID
				sd.GameName = "财神推金币"
				sd.DownBetTime = nowTime
				sd.StartTime = nowTime
				sd.EndTime = nowTime
				sd.PackageId = p.PackageId
				sd.WinStatementTotal = p.WinResultMoney
				sd.BetMoney = p.DownBet
				InsertStatementDB(sd)

				if p.PackageId != 11 && p.PackageId != 8 {
					// 插入盈余数据
					sur := &SurplusPoolDB{}
					sur.UpdateTime = time.Now()
					sur.TimeNow = time.Now().Format("2006-01-02 15:04:05")
					sur.Rid = p.RoomId
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
		}()
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
		// 换房清空累计金币
		p.ProgressBet = 0
		p.DownBetCount = 0
		// 保存区间节点位置
		room.ConfigPlace[room.Config] = m.Coordinates
		// 修改当前配置区间
		room.Config = m.ChangeCfg
		// 判断该金币区间是否存在金币位置存储，如果存在则返回，不存在则返回空
		change := &msg.ChangeRoomCfg_S2C{}
		if len(room.CoinList[room.Config]) == len(room.ConfigPlace[room.Config]) {
			var reset bool
			for _, v := range room.ConfigPlace[room.Config] {
				y, _ := strconv.ParseFloat(v.Location[2], 64)
				if y > 98 || y < -365 {
					reset = true
				}
			}
			if reset == true {
				var isHave bool
				for _, v := range room.CoinList[room.Config] {
					if v == FuDai {
						isHave = true
					}
				}

				room.CoinList[room.Config] = nil
				room.ConfigPlace[room.Config] = nil
				for i := 1; i <= 100; i++ {
					room.CoinNum[room.Config] ++
					room.CoinList[room.Config] = append(room.CoinList[room.Config], Coin+strconv.Itoa(int(room.CoinNum[room.Config])))
				}
				room.ConfigPlace[room.Config] = room.PushPlace

				if isHave == true {
					room.CoinList[room.Config] = append(room.CoinList[room.Config], FuDai)
					coinPlace := make([]*msg.Coordinate, 0)
					coinPlace = room.PushPlace
					place := &msg.Coordinate{}
					place.Location = []string{"18.807628842146926", "-107.9611168873588", "21"}
					coinPlace = append(coinPlace, place)
					room.ConfigPlace[room.Config] = coinPlace
				}
				change.IsChange = true
			}
		} else {
			var isHave bool
			for _, v := range room.CoinList[room.Config] {
				if v == FuDai {
					isHave = true
				}
			}
			room.CoinList[room.Config] = nil
			room.ConfigPlace[room.Config] = nil
			for i := 1; i <= 100; i++ {
				room.CoinNum[room.Config] ++
				room.CoinList[room.Config] = append(room.CoinList[room.Config], Coin+strconv.Itoa(int(room.CoinNum[room.Config])))
			}
			room.ConfigPlace[room.Config] = room.PushPlace

			if isHave == true {
				room.CoinList[room.Config] = append(room.CoinList[room.Config], FuDai)
				coinPlace := make([]*msg.Coordinate, 0)
				coinPlace = room.PushPlace
				place := &msg.Coordinate{}
				place.Location = []string{"18.807628842146926", "-107.9611168873588", "21"}
				coinPlace = append(coinPlace, place)
				room.ConfigPlace[room.Config] = coinPlace
			}
			change.IsChange = true
		}
		change.IsChange = true
		change.CoinList = room.CoinList[room.Config]
		change.Coordinates = room.ConfigPlace[room.Config]
		p.SendMsg(change)
	}
}

func (p *Player) SaveCoordinate(m *msg.SendCoordinate_C2S) {
	rid, _ := hall.UserRoom.Load(p.Id)
	v, _ := hall.RoomRecord.Load(rid)
	if v != nil {
		room := v.(*Room)
		// 保存区间节点位置
		room.ConfigPlace[room.Config] = m.Coordinates
		log.Debug("===>当前金币长度:%v,位置长度:%v", len(room.CoinList[room.Config]), len(room.ConfigPlace[room.Config]))
	}
}
