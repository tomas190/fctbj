package internal

import (
	"fmt"
	"strconv"
)

func test() {

	var totalBetWin float64
	var totalBetLose float64
	var totalWinNum int
	var totalLoseNum int

	var surplusPool float64 = 1000

	for i := 0; i < 10000; i++ {
		fmt.Println(i)
		p := &Player{}
		p.Init()
		r := &Room{}
		r.Config = "8" // 10金币
		for i := 1; i <= 100; i++ {
			r.CoinNum[r.Config] ++
			r.CoinList[r.Config] = append(r.CoinList[r.Config], Coin+strconv.Itoa(int(r.CoinNum["8"])))

		}
		var bet float64 = 10

		loseRate := 70
		percentageWin := 0
		countWin := 0
		percentageLose := 100
		countLose := 2

		num := RandInRange(1, 101)
		if num >= 0 { // 玩家赢钱
			settle := p.test2(r, bet)
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
						settle := p.test2(r, bet)
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
						settle := p.test2(r, bet)
						if settle > 0 { // 玩家赢
							// 盈余池判定
							if surplusPool > settle { // 盈余池足够
								break
							} else {                             // 盈余池不足
								if loseRateNum > int(loseRate) { // 30%玩家赢钱
									for {
										settle := p.test2(r, bet)
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

		taxWinMoney := totalBetWin - (totalBetWin * 0.05)
		fmt.Println("玩家总赢局:", totalWinNum)
		fmt.Println("玩家总输局:", totalLoseNum)
		fmt.Println("玩家总输:", int64(totalBetLose))
		fmt.Println("玩家总赢:", int64(totalBetWin))
		fmt.Println("玩家总赢(税后):", int64(taxWinMoney))
		fmt.Println("总流水:", int64(totalBetWin+totalBetLose))
	}
}

func (p *Player) test2(room *Room, bet float64) float64 {
	// 房间配置金额
	var goldNum int
	for { // 循环获取随机金币,避免随机到金币大于桌面金币数量
		goldNum = p.randGoldNum()
		if len(room.CoinList[room.Config]) >= goldNum {
			break
		}
	}
	p.DownBetList = room.CoinList[room.Config][:goldNum]
	//log.Debug("随机金币的数量:%v,切片值:%v", goldNum, p.DownBetList)
	settle := bet * float64(goldNum)
	return settle
}
