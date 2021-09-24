package internal

import (
	"fctbj/msg"
)

func GetRICH(money float64) (float64, float64) {
	num := RandInRange(1, 101)
	var rate float64
	if num >= 1 && num <= 70 {
		rate = 200
	} else if num >= 71 && num <= 85 {
		rate = 250
	} else if num >= 86 && num <= 95 {
		rate = 300
	} else if num >= 96 && num <= 98 {
		rate = 400
	} else if num >= 99 && num <= 100 {
		rate = 500
	}
	rich := rate * money
	return rate, rich
}

func (r *Room) GetPUSH(money float64) (float64, float64, int, int) {
	fudai1 := RandInRange(1, 4)
	fudai2 := RandInRange(0, 2)
	for i := 0; i < fudai1; i++ {
		r.CoinList[r.Config] = append(r.CoinList[r.Config], FuDai)
	}
	for i := 0; i < fudai2; i++ {
		r.CoinList[r.Config] = append(r.CoinList[r.Config], FuDai2)
	}

	var winNum int
	var luckyBag1 int
	var luckyBag2 int
	for _, v := range r.CoinList[r.Config] {
		if v == FuDai {
			luckyBag1++
		} else if v == FuDai2 {
			luckyBag2++
		}
	}
	for {
		winNum = RandInRange(50, 200)
		if winNum+(luckyBag1*LuckyBag)+(luckyBag2*LuckyBag2) <= 200 {
			break
		}
	}

	winMoney := (float64(winNum) * money) + (money * float64(luckyBag1*LuckyBag)) + (money * float64(luckyBag2*LuckyBag2))
	rate := winMoney / CfgMoney[r.Config]
	return rate, winMoney, fudai1, fudai2
}

func GetLUCKY(money float64) (float64, *msg.ThreePig) {
	num := RandInRange(1, 101)
	var rate float64
	if num >= 1 && num <= 50 {
		rate = 30
	} else if num >= 51 && num <= 80 {
		rate = 40
	} else if num >= 81 && num <= 100 {
		rate = 50
	}
	data := &msg.ThreePig{}
	if rate == 30 {
		data.PigSuccess = 30 * money
		data.PigLoser_1 = 40 * money
		data.PigLoser_2 = 50 * money
	} else if rate == 40 {
		data.PigSuccess = 40 * money
		data.PigLoser_1 = 30 * money
		data.PigLoser_2 = 50 * money
	} else if rate == 50 {
		data.PigSuccess = 50 * money
		data.PigLoser_1 = 30 * money
		data.PigLoser_2 = 40 * money
	}
	return rate, data
}

func GetGOLD(betNum int32) (int32, float64) {
	num := RandInRange(1, 101)
	var rate float64
	if num >= 1 && num <= 40 {
		rate = 2
	} else if num >= 41 && num <= 70 {
		rate = 3
	} else if num >= 71 && num <= 85 {
		rate = 4
	} else if num >= 86 && num <= 95 {
		rate = 5
	} else if num >= 96 && num <= 100 {
		rate = 6
	}
	money := rate * float64(betNum)
	return int32(rate), money
}
