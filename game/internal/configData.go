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
