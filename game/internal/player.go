package internal

import (
	"fctbj/conf"
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

}

func (p *Player) Init() {
	p.RoomId = ""
	p.RoundId = ""
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
	num := RandInRange(1, 10000)
	return fmt.Sprintf("%+v-%+v-%+v", time.Now().Unix(), num, p.Id)
}

func (p *Player) RespEnterRoom() {
	rid, _ := hall.UserRoom.Load(p.Id)
	v, _ := hall.RoomRecord.Load(rid)
	if v != nil {
		room := v.(*Room)
		p.IsExist = true
		enter := &msg.EnterRoom_S2C{}
		// 判断该金币区间是否存在金币位置存储，如果存在则返回，不存在则返回空
		if room.ConfigPlace[room.Config] != nil {
			var reset bool
			if len(room.CoinList[room.Config]) != len(room.ConfigPlace[room.Config]) {
				reset = true
			}
			for _, v := range room.ConfigPlace[room.Config] {
				y, _ := strconv.ParseFloat(v.Location[1], 64)
				if y > 98 || y < -365 {
					reset = true
				}
			}
			var isHave bool
			for _, v := range room.CoinList[room.Config] {
				if v == FuDai {
					isHave = true
				}
			}
			if reset == true {
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
					place.Location = []string{"82.02063531614465", "-235.30818435638344", "31"}
					coinPlace = append(coinPlace, place)
					room.ConfigPlace[room.Config] = coinPlace
				}
			}
			enter.IsChange = true
			enter.Coordinates = room.ConfigPlace[room.Config]
		} else {
			enter.IsChange = true
			room.ConfigPlace[room.Config] = room.PushPlace
			enter.Coordinates = room.ConfigPlace[room.Config]
			log.Debug("房间位置信息为空,添加预设值:%v", len(room.ConfigPlace[room.Config]))
		}
		log.Debug("返回房间数据,当前金币长度:%v,位置长度:%v", len(room.CoinList[room.Config]), len(room.ConfigPlace[room.Config]))

		enter.RoomData = room.RespRoomData()
		enter.IsPickGod = room.IsPickGod
		enter.IsLuckyPig = room.IsLuckyPig
		p.SendMsg(enter)
	}
}

//InsertPlayerData 插入玩家数据
func (p *Player) HandlePlayerData() {
	nowTime := time.Now().Unix()
	if p.TotalLoseMoney > 0 {
		p.RoundId = fmt.Sprintf("%+v-%+v", time.Now().Unix(), p.Id)
		loseReason := "发财推币机输钱"
		c2c.UserSyncLoseScore(p, nowTime, p.RoundId, loseReason, p.TotalLoseMoney)
	}

	if p.TotalWinMoney > 0 {
		p.RoundId = p.RandRoundId()
		winReason := "发财推币机赢钱"
		c2c.UserSyncWinScore(p, nowTime, p.RoundId, winReason, p.TotalWinMoney)
	}

	pac := packageTax[p.PackageId]
	taxR := pac / 100
	tax := p.TotalWinMoney * taxR
	resultMoney := (p.TotalWinMoney - tax) - p.TotalLoseMoney

	// 跑马灯
	if resultMoney > PaoMaDeng {
		c2c.NoticeWinMoreThan(p.Id, p.NickName, resultMoney)
	}

	// 插入运营数据
	data := &PlayerDownBetRecode{}
	data.Id = p.Id
	data.GameId = conf.Server.GameID
	data.RoundId = p.RoundId
	data.RoomId = p.RoomId
	data.DownBetInfo = p.DownBet
	data.DownBetTime = nowTime - 180
	data.StartTime = nowTime - 180
	data.EndTime = nowTime
	data.SettlementFunds = resultMoney
	data.SpareCash = p.Account
	data.TaxRate = taxR
	InsertAccessData(data)

	// 插入游戏统计数据
	sd := &StatementData{}
	sd.Id = p.Id
	sd.GameId = conf.Server.GameID
	sd.GameName = "财神推金币"
	sd.DownBetTime = nowTime - 180
	sd.StartTime = nowTime - 180
	sd.EndTime = nowTime
	sd.PackageId = p.PackageId
	sd.WinStatementTotal = p.TotalWinMoney
	sd.LoseStatementTotal = p.TotalLoseMoney
	sd.BetMoney = p.DownBet
	InsertStatementDB(sd)

	// 插入盈余数据
	sur := &SurplusPoolDB{}
	sur.UpdateTime = time.Now()
	sur.TimeNow = time.Now().Format("2006-01-02 15:04:05")
	sur.Rid = p.RoomId
	sur.PlayerNum = LoadPlayerCount()
	surPool := FindSurplusPool()
	if surPool != nil {
		sur.HistoryWin = surPool.HistoryWin
		sur.HistoryLose = surPool.HistoryLose
	}
	sur.HistoryWin += Decimal(p.TotalWinMoney)
	sur.TotalWinMoney += Decimal(p.TotalWinMoney)
	sur.HistoryLose += Decimal(p.TotalLoseMoney)
	sur.TotalLoseMoney += Decimal(p.TotalLoseMoney)
	InsertSurplusPool(sur)

	// 清除玩家累计数据
	p.DownBet = 0
	p.TotalWinMoney = 0
	p.TotalLoseMoney = 0
}
