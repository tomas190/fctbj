package internal

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	RECODE_PlayerInfoIDIsNull = "1001" // 玩家ID为空
	RECODE_UserMoneyNotEnough = "1002" // 玩家金额不足
	RECODE_RoomCfgMoneyERROR  = "1003" // 房间配置金额不对
	RECODE_TableNotHaveGold   = "1004" // 玩家桌面金币为空
	RECODE_ActionCoinNotHave  = "1005" // 玩家行动金币为空
)

var CfgMoney = map[string]float64{
	"1":  0.1,
	"2":  0.3,
	"3":  0.5,
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
