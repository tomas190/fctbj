package internal

import (
	"fctbj/msg"
	"fmt"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"strconv"
	"time"
)

type Player struct {
	// 玩家代理链接
	ConnAgent gate.Agent

	Id        string
	NickName  string
	HeadImg   string
	Account   float64
	Password  string
	Token     string
	RoomId    string
	RoundId   string
	RoundNum  int32
	PackageId uint16

	WinResultMoney  float64
	LoseResultMoney float64

	DownBet        float64  // 累计下注
	DownBetCount   int32    // 累计下注次数
	TotalWinMoney  float64  // 累计赢钱
	TotalLoseMoney float64  // 累计输钱
	ProgressBet    int32    // 掉落金币累计
	OffLineTime    int      // 离线时间
	IsExist        bool     // 是否存在房间内
	DownBetList    []string // 掉落金币切片
	IsLogin        bool     // 玩家是否登入
}

func (p *Player) Init() {
	p.RoomId = ""
	p.RoundId = ""
	p.RoundNum = 0
	p.WinResultMoney = 0
	p.LoseResultMoney = 0
	p.DownBet = 0
	p.DownBetCount = 0
	p.TotalWinMoney = 0
	p.TotalLoseMoney = 0
	p.ProgressBet = 0
	p.OffLineTime = -1
	p.IsExist = false
	p.DownBetList = nil
	p.IsLogin = true
}

//SendMsg 玩家向客户端发送消息
func (p *Player) SendMsg(msg interface{}) {
	if p.ConnAgent != nil {
		p.ConnAgent.WriteMsg(msg)
	}
}

//SendErrMsg 发送错误消息值
func (p *Player) SendErrMsg(errData string) {
	data := &msg.ErrorMsg_S2C{}
	data.MsgData = errData
	p.SendMsg(data)
}

func (p *Player) RandRoundId() string {
	p.RoundNum++
	return fmt.Sprintf("%+v-%+v-%+v", time.Now().Unix(), p.RoundNum, p.Id)
}

func (p *Player) RespEnterRoom() {
	rid, _ := hall.UserRoom.Load(p.Id)
	v, _ := hall.RoomRecord.Load(rid)
	if v != nil {
		room := v.(*Room)
		p.IsExist = true
		// 判断该金币区间是否存在金币位置存储，如果存在则返回，不存在则返回空
		enter := &msg.EnterRoom_S2C{}
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
				enter.IsChange = true
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
			enter.IsChange = true
		}
		log.Debug("返回房间数据,当前金币长度:%v,位置长度:%v", len(room.CoinList[room.Config]), len(room.ConfigPlace[room.Config]))
		enter.RoomData = room.RespRoomData()
		enter.Coordinates = room.ConfigPlace[room.Config]
		enter.IsPickGod = room.IsPickGod
		enter.IsLuckyPig = room.IsLuckyPig
		p.SendMsg(enter)
	}
}

// 获取玩家税收
func (p *Player) GetPackageTax() float64 {
	pac := packageTax[p.PackageId]
	taxR := pac / 100
	return taxR
}

// 玩家发送中心服赢钱
func (p *Player) SendC2CWinScore(winReason string, money float64) {
	p.RoundId = p.RandRoundId()
	c2c.UserSyncWinScore(p, GetTimeUnixNano(), p.RoundId, winReason, money)
}

// 玩家发送中心服输钱
func (p *Player) SendC2CLoseScore(money float64) {
	p.RoundId = p.RandRoundId()
	loseReason := "发财推币机输钱"
	c2c.UserSyncLoseScore(p, GetTimeUnixNano(), p.RoundId, loseReason, money)
}
