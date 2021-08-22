package internal

import (
	"fctbj/conf"
	"fctbj/msg"
	"fmt"
	"github.com/name5566/leaf/gate"
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

	DownBet        float64 // 累计下注
	DownBetCount   int32   // 累计下注次数
	TotalWinMoney  float64 // 累计赢钱
	TotalLoseMoney float64 // 累计输钱
	ProgressBet    int32   // 掉落金币累计
	OffLineTime    int     // 离线时间
	IsExist        bool    // 是否存在房间内

	DownBetList []string // 掉落金币切片

	ConfigPlace map[string][]*msg.Coordinate
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

	p.ConfigPlace = make(map[string][]*msg.Coordinate)
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

//InsertPlayerData 插入玩家数据
func (p *Player) HandlePlayerData() {
	nowTime := time.Now().Unix()
	if p.TotalLoseMoney > 0 {
		p.RoundId = fmt.Sprintf("%+v-%+v", time.Now().Unix(), p.Id)
		loseReason := "发财推币机输钱"
		c2c.UserSyncLoseScore(p, nowTime, p.RoundId, loseReason, p.TotalLoseMoney)
	}

	if p.TotalWinMoney > 0 {
		p.RoundId = fmt.Sprintf("%+v-%+v", time.Now().Unix(), p.Id)
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
