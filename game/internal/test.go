package internal

import (
	"fctbj/msg"
	"fmt"
	"github.com/name5566/leaf/log"
	"strconv"
)

//var surplusPool float64 = 0
//var surplusPool float64 = 1000000
var surplusPool float64 = -1000000

func test() {

	var totalBetWin float64
	var totalBetLose float64
	var totalWinNum int
	var totalLoseNum int

	var ProgressBet = 0
	var gameNumber = 0

	for i := 0; i < 100000; i++ {
		fmt.Println(i)
		p := &Player{}
		p.Init()
		p.DownBetList = nil

		r := &Room{}
		r.CoinNum = make(map[string]int32)
		r.CoinList = make(map[string][]string)

		r.Config = "8" // 10金币
		for i := 1; i <= 100; i++ {
			r.CoinNum[r.Config] ++
			r.CoinList[r.Config] = append(r.CoinList[r.Config], Coin+strconv.Itoa(int(r.CoinNum["8"])))

		}
		var bet float64 = 10
		ProgressBet++
		var result float64

		// 盈余设定
		p.test3(r, bet)

		result -= bet
		totalBetLose += bet
		surplusPool += bet

		if p.DownBetList != nil {
			win := float64(len(p.DownBetList)) * bet
			result += win
			totalBetWin += win
			surplusPool -= win
			fmt.Println("盈余池金额:", surplusPool)
		}

		if bet*Rate <= surplusPool {
			if ProgressBet >= 50 {
				ProgressBet = 0
				gameNumber++

				// 房间配置金额
				cfgMoney := CfgMoney[r.Config]
				data := &msg.GetRewards_S2C{}
				var winMoney float64
				num := RandInRange(0, 100)
				if num >= 0 && num <= 5 {
					data.RewardsNum = GOLD
				} else if num >= 6 && num <= 12 {
					data.RewardsNum = RICH
					_, data.GetMoney = GetRICH(cfgMoney)
					winMoney = data.GetMoney
				} else if num >= 13 && num <= 30 {
					data.RewardsNum = PUSH
					_, winMoney, _, _ = r.GetPUSH(cfgMoney)
				} else if num >= 31 && num <= 100 {
					data.RewardsNum = LUCKY
					_, data.LuckyPig = GetLUCKY(cfgMoney)
					winMoney = data.LuckyPig.PigSuccess
				}
				result += winMoney
				totalBetWin += winMoney
				surplusPool -= winMoney
				fmt.Println("盈余池金额3:", surplusPool)
			}
		}

		if result > 0 {
			totalWinNum++
		} else if result < 0 {
			totalLoseNum++
		}
	}

	fmt.Println("玩家小游戏局:", gameNumber)
	fmt.Println("玩家总赢局:", totalWinNum)
	fmt.Println("玩家总输局:", totalLoseNum)
	fmt.Println("玩家总输:", int64(totalBetLose))
	fmt.Println("玩家总赢:", int64(totalBetWin))
	fmt.Println("总流水:", int64(totalBetWin+totalBetLose))
}

func (p *Player) test3(r *Room, bet float64) {
	loseRate := 70

	percentageWin := 0.9 * 100
	countWin := 3
	percentageLose := 0 * 100
	countLose := 0

	num := RandInRange(1, 101)
	if num >= 50 { // 玩家赢钱
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
						break
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
}

func (p *Player) test2(room *Room, bet float64) float64 {
	// 房间配置金额
	var goldNum int
	for { // 循环获取随机金币,避免随机到金币大于桌面金币数量
		num := RandInRange(1, 101)
		if num >= 50 {
			goldNum = RandInRange(1, 10)
			if len(room.CoinList[room.Config]) >= goldNum {
				break
			}
		} else {
			p.DownBetList = nil
			break
		}
	}
	p.DownBetList = room.CoinList[room.Config][:goldNum]
	log.Debug("随机金币的数量:%v", goldNum)
	log.Debug("列表长度DownBetList:%v", p.DownBetList)
	settle := bet * float64(goldNum)
	return settle
}
