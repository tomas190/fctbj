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

