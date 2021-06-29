package internal

import (
	"C"
	"fctbj/conf"
	"fctbj/msg"
	"github.com/name5566/leaf/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var (
	session *mgo.Session
)

const (
	dbName          = "fctbj-Game"
	playerInfo      = "playerInfo"
	settleWinMoney  = "settleWinMoney"
	settleLoseMoney = "settleLoseMoney"
	surPlusDB       = "surPlusDB"
	surPool         = "surplus-pool"
)

// 连接数据库集合的函数 传入集合 默认连接IM数据库
func InitMongoDB() {
	// 此处连接正式线上数据库  下面是模拟的直接连接
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{conf.Server.MongoDBAddr},
		Timeout:  60 * time.Second,
		Database: conf.Server.MongoDBAuth,
		Username: conf.Server.MongoDBUser,
		Password: conf.Server.MongoDBPwd,
	}

	var err error
	session, err = mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatal("Connect DataBase 数据库连接ERROR: %v ", err)
	}
	log.Debug("Connect DataBase 数据库连接SUCCESS~")

	//打开数据库
	session.SetMode(mgo.Monotonic, true)
}

func connect(dbName, cName string) (*mgo.Session, *mgo.Collection) {
	s := session.Copy()
	c := s.DB(dbName).C(cName)
	return s, c
}

func (p *Player) FindPlayerInfo() {
	s, c := connect(dbName, playerInfo)
	defer s.Close()

	player := &msg.PlayerInfo{}
	player.Id = p.Id
	player.NickName = p.NickName
	player.HeadImg = p.HeadImg
	player.Account = p.Account

	err := c.Find(bson.M{"id": player.Id}).One(player)
	if err != nil {
		err2 := InsertPlayerInfo(player)
		if err2 != nil {
			log.Error("<----- 插入用户信息数据失败 ~ ----->:%v", err)
			return
		}
		log.Debug("<----- 插入用户信息数据成功 ~ ----->")
	}
}

func InsertPlayerInfo(player *msg.PlayerInfo) error {
	s, c := connect(dbName, playerInfo)
	defer s.Close()

	err := c.Insert(player)
	return err
}

//LoadPlayerCount 获取玩家数量
func LoadPlayerCount() int32 {
	s, c := connect(dbName, playerInfo)
	defer s.Close()

	n, err := c.Find(nil).Count()
	if err != nil {
		log.Debug("not Found Player Count, Maybe don't have Player")
		return 0
	}
	return int32(n)
}

//InsertWinMoney 插入赢钱数据
func InsertWinMoney(base interface{}) {
	s, c := connect(dbName, settleWinMoney)
	defer s.Close()

	err := c.Insert(base)
	if err != nil {
		log.Error("<----- 赢钱结算数据插入失败 ~ ----->:%v", err)
		return
	}
	log.Debug("<----- 赢钱结算数据插入成功 ~ ----->")
}

//InsertLoseMoney 插入输钱数据
func InsertLoseMoney(base interface{}) {
	s, c := connect(dbName, settleLoseMoney)
	defer s.Close()

	err := c.Insert(base)
	if err != nil {
		log.Error("<----- 输钱结算数据插入失败 ~ ----->:%v", err)
		return
	}
	log.Debug("<----- 输钱结算数据插入成功 ~ ----->")
}

//盈余池数据存入数据库
type SurplusPoolDB struct {
	UpdateTime     time.Time
	TimeNow        string  //记录时间（分为时间戳/字符串显示）
	Rid            string  //房间ID
	TotalWinMoney  float64 //玩家当局总赢
	TotalLoseMoney float64 //玩家当局总输
	PoolMoney      float64 //盈余池
	HistoryWin     float64 //玩家历史总赢
	HistoryLose    float64 //玩家历史总输
	PlayerNum      int32   //历史玩家人数
}

func FindSurplusPool() *SurplusPoolDB {
	s, c := connect(dbName, surPlusDB)
	defer s.Close()

	//c.RemoveAll(nil) // todo

	sur := &SurplusPoolDB{}
	err := c.Find(nil).Sort("-updatetime").One(sur)
	if err != nil {
		log.Error("<----- 查找SurplusPool数据失败 ~ ----->:%v", err)
		return nil
	}

	return sur
}

//InsertSurplusPool 插入盈余池数据
func InsertSurplusPool(sur *SurplusPoolDB) {
	s, c := connect(dbName, surPlusDB)
	defer s.Close()

	sur.PoolMoney = (sur.HistoryLose - (sur.HistoryWin * 1)) * 0.5
	log.Debug("surplusPoolDB 数据: %v", sur)
	err := c.Insert(sur)
	if err != nil {
		log.Error("<----- 数据库插入SurplusPool数据失败 ~ ----->:%v", err)
		return
	}
	log.Debug("<----- 数据库插入SurplusPool数据成功 ~ ----->")

	SurPool := &SurPool{}
	SurPool.GameId = conf.Server.GameID
	SurPool.SurplusPool = Decimal(sur.PoolMoney)
	SurPool.PlayerTotalLoseWin = Decimal(sur.HistoryLose - sur.HistoryWin)
	SurPool.PlayerTotalLose = Decimal(sur.HistoryLose)
	SurPool.PlayerTotalWin = Decimal(sur.HistoryWin)
	SurPool.TotalPlayer = sur.PlayerNum
	SurPool.FinalPercentage = 0.5
	SurPool.PercentageToTotalWin = 1
	SurPool.CoefficientToTotalPlayer = sur.PlayerNum * 0
	SurPool.PlayerLoseRateAfterSurplusPool = 0.7
	SurPool.DataCorrection = 0
	FindSurPool(SurPool)
}

type SurPool struct {
	GameId                         string  `json:"game_id" bson:"game_id"`
	PlayerTotalLose                float64 `json:"player_total_lose" bson:"player_total_lose"`
	PlayerTotalWin                 float64 `json:"player_total_win" bson:"player_total_win"`
	PercentageToTotalWin           float64 `json:"percentage_to_total_win" bson:"percentage_to_total_win"`
	TotalPlayer                    int32   `json:"total_player" bson:"total_player"`
	CoefficientToTotalPlayer       int32   `json:"coefficient_to_total_player" bson:"coefficient_to_total_player"`
	FinalPercentage                float64 `json:"final_percentage" bson:"final_percentage"`
	PlayerTotalLoseWin             float64 `json:"player_total_lose_win" bson:"player_total_lose_win" `
	SurplusPool                    float64 `json:"surplus_pool" bson:"surplus_pool"`
	PlayerLoseRateAfterSurplusPool float64 `json:"player_lose_rate_after_surplus_pool" bson:"player_lose_rate_after_surplus_pool"`
	DataCorrection                 float64 `json:"data_correction" bson:"data_correction"`
}

func FindSurPool(SurP *SurPool) {
	s, c := connect(dbName, surPool)
	defer s.Close()

	//c.RemoveAll(nil) // todo

	sur := &SurPool{}
	err := c.Find(nil).One(sur)
	if err != nil {
		InsertSurPool(SurP)
	} else {
		SurP.SurplusPool = (SurP.PlayerTotalLose - (SurP.PlayerTotalWin * sur.PercentageToTotalWin) - float64(SurP.TotalPlayer*sur.CoefficientToTotalPlayer) + sur.DataCorrection) * sur.FinalPercentage
		SurP.FinalPercentage = sur.FinalPercentage
		SurP.PercentageToTotalWin = sur.PercentageToTotalWin
		SurP.CoefficientToTotalPlayer = sur.CoefficientToTotalPlayer
		SurP.PlayerLoseRateAfterSurplusPool = sur.PlayerLoseRateAfterSurplusPool
		SurP.DataCorrection = sur.DataCorrection
		UpdateSurPool(SurP)
	}
}

//插入盈余池统一字段
func InsertSurPool(sur *SurPool) {
	s, c := connect(dbName, surPool)
	defer s.Close()

	log.Debug("SurPool 数据: %v", sur)

	err := c.Insert(sur)
	if err != nil {
		log.Error("<----- 数据库插入SurPool数据失败 ~ ----->:%v", err)
		return
	}
	log.Debug("<----- 数据库插入SurPool数据成功 ~ ----->")
}

func UpdateSurPool(sur *SurPool) {
	s, c := connect(dbName, surPool)
	defer s.Close()

	err := c.Update(bson.M{}, sur)
	if err != nil {
		log.Error("<----- 更新 SurPool数据失败 ~ ----->:%v", err)
		return
	}
	log.Debug("<----- 更新SurPool数据成功 ~ ----->")
}

func GetSurPlusMoney() float64 {
	s, c := connect(dbName, surPool)
	defer s.Close()

	//c.RemoveAll(nil) // todo

	sur := &SurPool{}
	err := c.Find(nil).One(sur)
	if err != nil {
		log.Debug("获取GetSurP数据失败:%v", err)
		return 0
	}
	return sur.SurplusPool
}