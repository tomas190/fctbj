package internal

import (
	"fctbj/msg"
	"math/rand"
	"time"
)

func GetRICH(money float64) (float64, float64) {
	slice := []float64{200, 200, 200, 200, 200, 200, 250, 250, 300, 500}
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(slice))
	rich := slice[n] * money
	return slice[n], rich
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
	data := &msg.ThreePig{}
	slice := []float64{30, 30, 30, 30, 30, 40, 40, 40, 50, 50}
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(slice))
	if slice[n] == 30 {
		data.PigSuccess = 30 * money
		data.PigLoser_1 = 40 * money
		data.PigLoser_2 = 50 * money
	} else if slice[n] == 40 {
		data.PigSuccess = 40 * money
		data.PigLoser_1 = 30 * money
		data.PigLoser_2 = 50 * money
	} else if slice[n] == 50 {
		data.PigSuccess = 50 * money
		data.PigLoser_1 = 30 * money
		data.PigLoser_2 = 40 * money
	}
	return slice[n], data
}

func GetGOLD(betNum int32) (int32, float64) {
	num := RandInRange(2, 7)
	push := float64(num) * float64(betNum)
	return int32(num), push
}
