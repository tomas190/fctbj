package internal

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	RECODE_CHAOCHUXIANHONG = "4444"
)

var CfgMoney = map[string]float64{
	"1":  0.1,
	"2":  0.3,
	"3":  0.4,
	"4":  1.0,
	"5":  3.0,
	"6":  5.0,
	"7":  8.0,
	"8":  10.0,
	"9":  20.0,
	"10": 30.0,
	"11": 40.0,
	"12": 50.0,
	"13": 60.0,
	"14": 80.0,
	"15": 100.0,
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.6f", value), 64)
	return value
}

func RandInRange(min int, max int) int {
	time.Sleep(1 * time.Nanosecond)
	return rand.Intn(max-min) + min
}

func SetPackageTaxM(packageT uint16, tax float64) {
	packageTax[packageT] = tax
}
