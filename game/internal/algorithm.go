package internal

import (
	"fctbj/msg"
	"math/rand"
	"time"
)

func GetRICH(money float64) float64 {
	slice := []float64{200, 200, 200, 200, 200, 200, 250, 250, 300, 500}
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(slice))
	rich := slice[n] * money
	return rich
}

func GetPUSH(money float64) float64 {
	num := RandInRange(70, 100)
	push := float64(num) * money
	return push
}

func GetLUCKY(money float64) *msg.ThreePig {
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
	return data
}

func GetGOLD(betNum int32) (int32, float64) {
	num := RandInRange(2, 7)
	push := float64(num) * float64(betNum)
	return int32(num), push
}
