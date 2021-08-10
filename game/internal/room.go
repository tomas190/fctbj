package internal

import (
	"fctbj/msg"
	"fmt"
	"math/rand"
	"time"
)

const (
	GOLD  = 1 // 接金币
	RICH  = 2 // 吐钱
	PUSH  = 3 // 财神推金币
	LUCKY = 4 // 三只小猪
)

const (
	Rate = 500 // 最高500倍率
)

const (
	Coin   = "coin"
	FuDai  = "fudai"
	FuDai2 = "fudai2"
)

const (
	LuckyBag  = 10
	LuckyBag2 = 20
	PaoZhu    = 30
	YuXi      = 40
	ShuiJing  = 50
)

var (
	packageTax map[uint16]float64
)

type Room struct {
	RoomId   string              // 房间号
	Config   string              // 房间配置
	Player   *Player             // 玩家信息
	CoinNum  map[string]int32    // coin序号
	CoinList map[string][]string // 金币列表
}

func (r *Room) Init() {
	r.RoomId = fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	r.Config = "1"
	r.CoinInit()
}

func (r *Room) CoinInit() {
	r.CoinNum["1"] = 0
	r.CoinNum["2"] = 0
	r.CoinNum["3"] = 0
	r.CoinNum["4"] = 0
	r.CoinNum["5"] = 0
	r.CoinNum["6"] = 0
	r.CoinNum["7"] = 0
	r.CoinNum["8"] = 0
	r.CoinNum["9"] = 0
	r.CoinNum["10"] = 0
	r.CoinNum["11"] = 0
	r.CoinNum["12"] = 0
	r.CoinNum["13"] = 0
	r.CoinNum["14"] = 0
	r.CoinNum["15"] = 0
	for i := 1; i <= 100; i++ {
		r.CoinNum["1"] += int32(i)
		r.CoinNum["2"] += int32(i)
		r.CoinNum["3"] += int32(i)
		r.CoinNum["4"] += int32(i)
		r.CoinNum["5"] += int32(i)
		r.CoinNum["6"] += int32(i)
		r.CoinNum["7"] += int32(i)
		r.CoinNum["8"] += int32(i)
		r.CoinNum["9"] += int32(i)
		r.CoinNum["10"] += int32(i)
		r.CoinNum["11"] += int32(i)
		r.CoinNum["12"] += int32(i)
		r.CoinNum["13"] += int32(i)
		r.CoinNum["14"] += int32(i)
		r.CoinNum["15"] += int32(i)
		r.CoinList["1"] = append(r.CoinList["1"], Coin+string(r.CoinNum["1"]))
		r.CoinList["2"] = append(r.CoinList["2"], Coin+string(r.CoinNum["2"]))
		r.CoinList["3"] = append(r.CoinList["3"], Coin+string(r.CoinNum["3"]))
		r.CoinList["4"] = append(r.CoinList["4"], Coin+string(r.CoinNum["4"]))
		r.CoinList["5"] = append(r.CoinList["5"], Coin+string(r.CoinNum["5"]))
		r.CoinList["6"] = append(r.CoinList["6"], Coin+string(r.CoinNum["6"]))
		r.CoinList["7"] = append(r.CoinList["7"], Coin+string(r.CoinNum["7"]))
		r.CoinList["8"] = append(r.CoinList["8"], Coin+string(r.CoinNum["8"]))
		r.CoinList["9"] = append(r.CoinList["9"], Coin+string(r.CoinNum["9"]))
		r.CoinList["10"] = append(r.CoinList["10"], Coin+string(r.CoinNum["10"]))
		r.CoinList["11"] = append(r.CoinList["11"], Coin+string(r.CoinNum["11"]))
		r.CoinList["12"] = append(r.CoinList["12"], Coin+string(r.CoinNum["12"]))
		r.CoinList["13"] = append(r.CoinList["13"], Coin+string(r.CoinNum["13"]))
		r.CoinList["14"] = append(r.CoinList["14"], Coin+string(r.CoinNum["14"]))
		r.CoinList["15"] = append(r.CoinList["15"], Coin+string(r.CoinNum["15"]))
	}
	r.CoinList["1"] = append(r.CoinList["1"], FuDai)
	r.CoinList["2"] = append(r.CoinList["2"], FuDai)
	r.CoinList["3"] = append(r.CoinList["3"], FuDai)
	r.CoinList["4"] = append(r.CoinList["4"], FuDai)
	r.CoinList["5"] = append(r.CoinList["5"], FuDai)
	r.CoinList["6"] = append(r.CoinList["6"], FuDai)
	r.CoinList["7"] = append(r.CoinList["7"], FuDai)
	r.CoinList["8"] = append(r.CoinList["8"], FuDai)
	r.CoinList["9"] = append(r.CoinList["9"], FuDai)
	r.CoinList["10"] = append(r.CoinList["10"], FuDai)
	r.CoinList["11"] = append(r.CoinList["11"], FuDai)
	r.CoinList["12"] = append(r.CoinList["12"], FuDai)
	r.CoinList["13"] = append(r.CoinList["13"], FuDai)
	r.CoinList["14"] = append(r.CoinList["14"], FuDai)
	r.CoinList["15"] = append(r.CoinList["15"], FuDai)
}

//RespRoomData 返回房间数据
func (r *Room) RespRoomData() *msg.RoomData {
	rd := &msg.RoomData{}
	rd.RoomId = r.RoomId
	rd.CfgId = r.Config
	rd.CoinList = r.CoinList[r.Config]
	//rd.PlayerInfo = new(msg.PlayerInfo)
	//rd.PlayerInfo.Id = r.Player.Id
	//rd.PlayerInfo.Account = r.Player.Account
	//rd.PlayerInfo.NickName = r.Player.NickName
	//rd.PlayerInfo.HeadImg = r.Player.HeadImg
	return rd
}

func (r *Room) ExistLuckyBag() bool {
	for _, v := range r.CoinList[r.Config] {
		if v == FuDai {
			return true
		}
	}
	return false
}
