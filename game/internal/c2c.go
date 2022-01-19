package internal

import (
	"encoding/json"
	"fctbj/conf"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/name5566/leaf/log"
	"gopkg.in/mgo.v2/bson"
)

// 防止并发写Websocket用的锁
var syncWrite sync.Mutex

//CGTokenRsp 接受Token结构体
type CGTokenRsp struct {
	Token string
}

//CGCenterRsp 中心返回消息结构体
type CGCenterRsp struct {
	Status string
	Code   int
	Msg    *CGTokenRsp
}

//Conn4Center 连接到Center(中心服务器)的网络协议处理器
type Conn4Center struct {
	GameId    string
	centerUrl string
	token     string
	DevKey    string
	conn      *websocket.Conn

	//除于登录成功状态
	LoginStat bool

	closebreathchan  chan bool
	closereceivechan chan bool

	//待处理的用户登录请求
	waitUser map[string]*UserCallback
}

//Init 初始化
func (c4c *Conn4Center) Init() {
	c4c.GameId = conf.Server.GameID
	c4c.DevKey = conf.Server.DevKey
	c4c.LoginStat = false

	c4c.waitUser = make(map[string]*UserCallback)
}

//CreatConnect 和Center建立链接
func (c4c *Conn4Center) CreatConnect() {
	c4c.centerUrl = conf.Server.CenterUrl

	log.Debug("--- dial: --- : %v", c4c.centerUrl)
	for {
		conn, rsp, err := websocket.DefaultDialer.Dial(c4c.centerUrl, nil)
		log.Debug("<--- Dial rsp --->: %v", rsp)
		if err == nil {
			c4c.conn = conn
			break
		}
		time.Sleep(time.Second * 5)
	}

	c4c.ServerLoginCenter()

	c4c.Run()
}

func (c4c *Conn4Center) ReConnect() {
	if c4c.LoginStat == true {
		return
	}
	time.Sleep(time.Second * 5)

	c4c.centerUrl = conf.Server.CenterUrl

	log.Debug("--- dial: --- : %v", c4c.centerUrl)
	for {
		conn, rsp, err := websocket.DefaultDialer.Dial(c4c.centerUrl, nil)
		log.Debug("<--- Dial rsp --->: %v", rsp)
		if err == nil {
			c4c.conn = conn
			break
		}
		time.Sleep(time.Second * 5)
	}

	c4c.ServerLoginCenter()
}

//Run 开始运行,监听中心服务器的返回
func (c4c *Conn4Center) Run() {
	ticker := time.NewTicker(time.Second * 3)
	log.Debug("发送心跳!")
	go func() {
		for { //循环
			select {
			case <-ticker.C:
				c4c.onBreath()
			case <-c4c.closebreathchan:
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case <-c4c.closereceivechan:
				return
			default:
				typeId, message, err := c4c.conn.ReadMessage()
				if err != nil {
					log.Debug("Here is error by ReadMessage~ %v", err)
					log.Error(err.Error())
				}
				if typeId == -1 {
					log.Debug("中心服异常消息~")
					c4c.LoginStat = false
					c4c.ReConnect()
				} else {
					c4c.onReceive(typeId, message)
				}
			}
		}
	}()
}

//onBreath 中心服心跳
func (c4c *Conn4Center) onBreath() {
	syncWrite.Lock()
	err := c4c.conn.WriteMessage(websocket.TextMessage, []byte(""))
	if err != nil {
		log.Error(err.Error())
	}
	syncWrite.Unlock()
}

//onReceive 接收消息
func (c4c *Conn4Center) onReceive(messType int, messBody []byte) {
	if messType == websocket.TextMessage {
		baseData := &BaseMessage{}

		decoder := json.NewDecoder(strings.NewReader(string(messBody)))
		decoder.UseNumber()

		err := decoder.Decode(&baseData)
		if err != nil {
			log.Error(err.Error())
		}

		switch baseData.Event {
		case msgServerLogin:
			c4c.onServerLogin(baseData.Data)
			break
		case msgUserLogin:
			c4c.onUserLogin(baseData.Data)
			break
		case msgUserLogout:
			c4c.onUserLogout(baseData.Data)
			break
		case msgUserWinScore:
			c4c.onUserWinScore(baseData.Data)
			break
		case msgUserLoseScore:
			c4c.onUserLoseScore(baseData.Data)
			break
		case msgLockSettlement:
			c4c.onLockSettlement(baseData.Data)
			break
		case msgUnlockSettlement:
			c4c.onUnlockSettlement(baseData.Data)
			break
		case msgWinMoreThanNotice:
			c4c.onWinMoreThanNotice(baseData.Data)
			break
		default:
			log.Error("Receive a message but don't identify~")
		}
	}
}

//onServerLogin 服务器登录
func (c4c *Conn4Center) onServerLogin(msgBody interface{}) {
	log.Debug("<-------- onServerLogin -------->")
	data, ok := msgBody.(map[string]interface{})
	if !ok {
		log.Debug("onServerLogin Error")
	}

	code, err := data["code"].(json.Number).Int64()
	if err != nil {
		log.Error(err.Error())
	}

	if data["status"] == "SUCCESS" && code == 200 {
		log.Debug("<-------- serverLogin SUCCESS~!!! -------->")
		c4c.LoginStat = true

		SendTgMessage("启动成功")

		msgInfo := data["msg"].(map[string]interface{})

		globals := msgInfo["globals"].([]interface{})
		for _, v := range globals {
			info := v.(map[string]interface{})

			var nPackage uint16
			var nTax float64

			jsonPackageId, err := info["package_id"].(json.Number).Int64()
			if err != nil {
				log.Debug("onServerLogin: jsonPackageId:%v", err.Error())
			} else {
				nPackage = uint16(jsonPackageId)
			}
			jsonTax, err := info["platform_tax_percent"].(json.Number).Float64()

			if err != nil {
				log.Debug("onServerLogin: jsonTax:%v", err.Error())
			} else {
				nTax = jsonTax
			}

			SetPackageTaxM(nPackage, nTax)
		}
	}

}

//onUserLogin 收到中心服的用户登录回应
func (c4c *Conn4Center) onUserLogin(msgBody interface{}) {
	data, ok := msgBody.(map[string]interface{})
	if !ok {
		log.Debug("onUserLogout Error: %v", data)
	}

	code, err := data["code"].(json.Number).Int64()
	if err != nil {
		log.Error(err.Error())
	}

	if code != 200 {
		log.Debug("同步中心服登录失败:%v", data)
		return
	}

	if data["status"] == "SUCCESS" && code == 200 {
		log.Debug("<-------- UserLogin SUCCESS~ -------->:%v", data)

		userInfo, ok := data["msg"].(map[string]interface{})
		var strId string
		var userData *UserCallback
		if ok {
			gameUser, uok := userInfo["game_user"].(map[string]interface{})
			if uok {
				nick := gameUser["game_nick"]
				headImg := gameUser["game_img"]
				userId := gameUser["id"]
				packageId := gameUser["package_id"]

				intID, err := userId.(json.Number).Int64()
				if err != nil {
					log.Fatal("onUserLogin intID:%v", err.Error())
				}
				strId = strconv.Itoa(int(intID))

				user, _ := hall.UserRecord.Load(strId)
				if user != nil {
					u := user.(*Player)
					u.IsLogin = true
				}

				pckId, err2 := packageId.(json.Number).Int64()
				if err2 != nil {
					log.Fatal("onUserLogin pckId:%v", err2.Error())
				}

				//找到等待登录玩家
				userData, ok = c4c.waitUser[strId]
				if ok {
					userData.Data.HeadImg = headImg.(string)
					userData.Data.NickName = nick.(string)
					userData.Data.PackageId = uint16(pckId)
				}
			}
			gameAccount, okA := userInfo["game_account"].(map[string]interface{})

			if okA {
				balance := gameAccount["balance"]
				floatBalance, err := balance.(json.Number).Float64()
				if err != nil {
					log.Error(err.Error())
				}

				userData.Data.Account = floatBalance

				//调用玩家绑定回调函数
				if userData.Callback != nil {
					userData.Callback(&userData.Data)
				}
			}
		}
	}
}

func (c4c *Conn4Center) onUserLogout(msgBody interface{}) {
	data, ok := msgBody.(map[string]interface{})
	if !ok {
		log.Debug("onUserLogout Error: %v", data)
	}

	code, err := data["code"].(json.Number).Int64()
	if err != nil {
		log.Error(err.Error())
	}

	if code != 200 {
		log.Debug("同步中心服登出失败:%v", data)
		return
	}

	if data["status"] == "SUCCESS" && code == 200 {
		log.Debug("<-------- UserLogout SUCCESS~ -------->:%v", data)

		userInfo, ok := data["msg"].(map[string]interface{})
		var strId string
		var userData *UserCallback
		if ok {
			gameUser, uok := userInfo["game_user"].(map[string]interface{})
			if uok {
				nick := gameUser["game_nick"]
				headImg := gameUser["game_img"]
				userId := gameUser["id"]

				intID, err := userId.(json.Number).Int64()
				if err != nil {
					log.Fatal("onUserLogout:%v", err.Error())
				}
				strId = strconv.Itoa(int(intID))

				user, _ := hall.UserRecord.Load(strId)
				if user != nil {
					u := user.(*Player)
					u.IsLogin = false
				}

				//找到等待登录玩家
				userData, ok = c4c.waitUser[strId]
				if ok {
					userData.Data.HeadImg = headImg.(string)
					userData.Data.NickName = nick.(string)
				}
			}
		}
	}
}

func (c4c *Conn4Center) onUserWinScore(msgBody interface{}) {
	data, ok := msgBody.(map[string]interface{})
	if !ok {
		log.Debug("onUserWinScore Error")
	}

	code, err := data["code"].(json.Number).Int64()
	if err != nil {
		log.Error(err.Error())
	}

	if code != 200 {
		log.Debug("同步中心服赢钱失败:%v", data)
		return
	}

	if data["status"] == "SUCCESS" && code == 200 {
		log.Debug("<-------- UserWinScore SUCCESS~ -------->")

		//将Win数据插入数据
		InsertWinMoney(msgBody) //todo

		userInfo, ok := data["msg"].(map[string]interface{})
		if ok {
			jsonScore := userInfo["final_pay"]
			score, err := jsonScore.(json.Number).Float64()

			log.Debug("同步中心服赢钱成功:%v", score)

			if err != nil {
				log.Error(err.Error())
				return
			}
		}
	}
}

func (c4c *Conn4Center) onUserLoseScore(msgBody interface{}) {
	data, ok := msgBody.(map[string]interface{})
	if !ok {
		log.Debug("onUserLoseScore Error")
	}

	code, err := data["code"].(json.Number).Int64()
	if err != nil {
		log.Error(err.Error())
	}
	msgData, ok := data["msg"].(map[string]interface{})
	if ok {
		order := msgData["order"]
		if code != 200 {
			log.Debug("同步中心服输钱失败:%v", data)
			v, ok := hall.OrderIDRecord.Load(order)
			if ok {
				p := v.(*Player)
				rid, _ := hall.UserRoom.Load(p.Id)
				v, _ := hall.RoomRecord.Load(rid)
				if v != nil {
					room := v.(*Room)
					message := fmt.Sprintf("玩家" + p.Id + "输钱失败并登出")
					SendTgMessage(message)
					p.ExitFromRoom(room)
					hall.OrderIDRecord.Delete(order)
					p.LoseChan <- false
				}
			}
			return
		}
		if data["status"] == "SUCCESS" && code == 200 {
			log.Debug("<-------- UserLoseScore SUCCESS~ -------->")

			v, ok := hall.OrderIDRecord.Load(order)
			if ok {
				p := v.(*Player)
				hall.OrderIDRecord.Delete(order)
				p.LoseChan <- true
			}

			//将Lose数据插入数据
			InsertLoseMoney(msgBody)

			userInfo, ok := data["msg"].(map[string]interface{})
			if ok {
				jsonScore := userInfo["final_pay"]
				score, err := jsonScore.(json.Number).Float64()

				log.Debug("同步中心服输钱成功:%v", score)

				if err != nil {
					log.Error(err.Error())
					return
				}
			}
		}
	}
}

//onWinMoreThanNotice 加锁金额
func (c4c *Conn4Center) onLockSettlement(msgBody interface{}) {
	data, ok := msgBody.(map[string]interface{})
	if ok {
		code, err := data["code"].(json.Number).Int64()
		if err != nil {
			log.Fatal("onLockSettlement:%v", err.Error())
		}

		fmt.Println(code, reflect.TypeOf(code))
		if data["status"] == "SUCCESS" && code == 200 {
			log.Debug("<-------- onLockSettlement SUCCESS~!!! -------->")
		}
	}
}

//onWinMoreThanNotice 解锁金额
func (c4c *Conn4Center) onUnlockSettlement(msgBody interface{}) {
	data, ok := msgBody.(map[string]interface{})
	if ok {
		code, err := data["code"].(json.Number).Int64()
		if err != nil {
			log.Fatal("onUnlockSettlement:%v", err.Error())
		}

		fmt.Println(code, reflect.TypeOf(code))
		if data["status"] == "SUCCESS" && code == 200 {
			log.Debug("<-------- onUnlockSettlement SUCCESS~!!! -------->")
		}
	}
}

func (c4c *Conn4Center) onBankerStatus(msgBody interface{}) {
	data, ok := msgBody.(map[string]interface{})
	if !ok {
		log.Debug("onBankerStatus Error")
	}

	code, err := data["code"].(json.Number).Int64()
	if err != nil {
		log.Error(err.Error())
	}

	if code != 200 {
		log.Debug("同步中心服庄家状态失败:%v", data)
		return
	}

	if data["status"] == "SUCCESS" && code == 200 {
		log.Debug("<-------- onBankerStatus SUCCESS~ -------->")
	}
}

func (c4c *Conn4Center) onBankerWinScore(msgBody interface{}) {
	data, ok := msgBody.(map[string]interface{})
	if !ok {
		log.Debug("onBankerWinScore Error")
	}

	code, err := data["code"].(json.Number).Int64()
	if err != nil {
		log.Error(err.Error())
	}

	if code != 200 {
		log.Debug("同步中心服庄家赢钱失败:%v", data)
		return
	}

	if data["status"] == "SUCCESS" && code == 200 {
		log.Debug("<-------- onBankerWinScore SUCCESS~ -------->")

		userInfo, ok := data["msg"].(map[string]interface{})
		if ok {
			jsonScore := userInfo["final_pay"]
			score, err := jsonScore.(json.Number).Float64()

			log.Debug("同步中心服庄家赢钱成功:%v", score)

			if err != nil {
				log.Error(err.Error())
				return
			}
		}
	}
}

func (c4c *Conn4Center) onBankerLoseScore(msgBody interface{}) {
	data, ok := msgBody.(map[string]interface{})
	if !ok {
		log.Debug("onBankerLoseScore Error")
	}

	code, err := data["code"].(json.Number).Int64()
	if err != nil {
		log.Error(err.Error())
	}
	if code != 200 {
		log.Error("同步中心服庄家输钱失败:%v", data)
		return
	}

	if data["status"] == "SUCCESS" && code == 200 {
		log.Debug("<-------- onBankerLoseScore SUCCESS~ -------->")

		userInfo, ok := data["msg"].(map[string]interface{})
		if ok {
			jsonScore := userInfo["final_pay"]
			score, err := jsonScore.(json.Number).Float64()

			log.Debug("同步中心服庄家输钱成功:%v", score)

			if err != nil {
				log.Error(err.Error())
				return
			}
		}
	}
}

//onWinMoreThanNotice 服务器登录
func (c4c *Conn4Center) onWinMoreThanNotice(msgBody interface{}) {
	data, ok := msgBody.(map[string]interface{})
	if ok {
		code, err := data["code"].(json.Number).Int64()
		if err != nil {
			log.Fatal("onWinMoreThanNotice:%v", err.Error())
		}

		fmt.Println(code, reflect.TypeOf(code))
		if data["status"] == "SUCCESS" && code == 200 {
			log.Debug("<-------- onWinMoreThanNotice SUCCESS~!!! -------->")
		}
	}
}

//SendMsg2Center 发送消息到中心服
func (c4c *Conn4Center) SendMsg2Center(data interface{}) {
	syncWrite.Lock()
	defer syncWrite.Unlock()
	// Json序列化
	codeData, err1 := json.Marshal(data)
	if err1 != nil {
		log.Error(err1.Error())
	}
	log.Debug("Msg to Send Center:%v", string(codeData))

	err2 := c4c.conn.WriteMessage(websocket.TextMessage, codeData)
	if err2 != nil {
		log.Fatal("SendMsg2Center:%v", err2.Error())
	}
}

//ServerLoginCenter 服务器登录Center
func (c4c *Conn4Center) ServerLoginCenter() {
	baseData := &BaseMessage{}
	baseData.Event = msgServerLogin
	baseData.Data = ServerLogin{
		Host:    conf.Server.CenterServer,
		Port:    conf.Server.CenterServerPort,
		GameId:  c4c.GameId,
		DevName: conf.Server.DevName,
		DevKey:  c4c.DevKey,
	}
	// 发送消息到中心服
	c4c.SendMsg2Center(baseData)
}

//UserLoginCenter 用户登录
func (c4c *Conn4Center) UserLoginCenter(userId string, password string, token string, callback func(data *Player)) {
	if !c4c.LoginStat {
		log.Debug("<-------- fctbj not ready~!!! -------->")
		return
	}
	baseData := &BaseMessage{}
	baseData.Event = msgUserLogin
	id, _ := strconv.Atoi(userId)
	if password != "" {
		baseData.Data = &UserReq{
			ID:       id,
			PassWord: password,
			GameId:   c4c.GameId,
			DevName:  conf.Server.DevName,
			DevKey:   c4c.DevKey}
	} else {
		baseData.Data = &UserReq{
			ID:      id,
			Token:   token,
			GameId:  c4c.GameId,
			DevName: conf.Server.DevName,
			DevKey:  c4c.DevKey}
	}

	c4c.SendMsg2Center(baseData)

	//加入待处理map，等待处理
	c4c.waitUser[userId] = &UserCallback{}
	c4c.waitUser[userId].Data.Id = userId
	c4c.waitUser[userId].Callback = callback
}

//UserLogoutCenter 用户登出
func (c4c *Conn4Center) UserLogoutCenter(userId string, password string, token string) {
	base := &BaseMessage{}
	base.Event = msgUserLogout
	id, _ := strconv.Atoi(userId)
	if password != "" {
		base.Data = &UserReq{
			ID:       id,
			PassWord: password,
			GameId:   c4c.GameId,
			DevName:  conf.Server.DevName,
			DevKey:   c4c.DevKey}
	} else {
		base.Data = &UserReq{
			ID:      id,
			Token:   token,
			GameId:  c4c.GameId,
			DevName: conf.Server.DevName,
			DevKey:  c4c.DevKey}
	}

	// 发送消息到中心服
	c4c.SendMsg2Center(base)
}

//UserSyncWinScore 同步赢分
func (c4c *Conn4Center) UserSyncWinScore(p *Player, timeUnix int64, roundId, reason string, betMoney float64) {
	baseData := &BaseMessage{}
	baseData.Event = msgUserWinScore
	id, _ := strconv.Atoi(p.Id)
	userWin := &UserChangeScore{}
	userWin.Auth.DevName = conf.Server.DevName
	userWin.Auth.DevKey = c4c.DevKey
	userWin.Info.CreateTime = timeUnix
	userWin.Info.GameId = c4c.GameId
	userWin.Info.ID = id
	userWin.Info.LockMoney = 0
	userWin.Info.Money = betMoney
	userWin.Info.BetMoney = 0
	userWin.Info.Order = bson.NewObjectId().Hex()

	userWin.Info.PayReason = reason
	userWin.Info.PreMoney = 0
	userWin.Info.RoundId = roundId
	baseData.Data = userWin
	c4c.SendMsg2Center(baseData)
}

//UserSyncWinScore 同步输分
func (c4c *Conn4Center) UserSyncLoseScore(p *Player, timeUnix int64, roundId, reason string, betMoney float64) {
	baseData := &BaseMessage{}
	baseData.Event = msgUserLoseScore
	id, _ := strconv.Atoi(p.Id)
	userLose := &UserChangeScore{}
	userLose.Auth.DevName = conf.Server.DevName
	userLose.Auth.DevKey = c4c.DevKey
	userLose.Info.CreateTime = timeUnix
	userLose.Info.GameId = c4c.GameId
	userLose.Info.ID = id
	userLose.Info.LockMoney = 0
	userLose.Info.Money = betMoney
	userLose.Info.BetMoney = betMoney
	userLose.Info.Order = bson.NewObjectId().Hex()
	userLose.Info.PayReason = reason
	userLose.Info.PreMoney = 0
	userLose.Info.RoundId = roundId
	baseData.Data = userLose
	c4c.SendMsg2Center(baseData)
	hall.OrderIDRecord.Store(userLose.Info.Order, p)
}

//锁钱
func (c4c *Conn4Center) LockSettlement(p *Player) {
	timeStr := time.Now().Format("2006-01-02_15:04:05")
	loseOrder := p.Id + "_" + timeStr + "_LockMoney"

	baseData := &BaseMessage{}
	baseData.Event = msgLockSettlement
	id, _ := strconv.Atoi(p.Id)
	lockMoney := &UserChangeScore{}
	lockMoney.Auth.DevName = conf.Server.DevName
	lockMoney.Auth.DevKey = c4c.DevKey
	lockMoney.Info.CreateTime = time.Now().Unix()
	lockMoney.Info.GameId = c4c.GameId
	lockMoney.Info.ID = id
	lockMoney.Info.LockMoney = 0
	lockMoney.Info.Money = 0
	lockMoney.Info.Order = loseOrder
	lockMoney.Info.PayReason = "lockMoney"
	lockMoney.Info.PreMoney = 0
	lockMoney.Info.RoundId = p.RoundId
	baseData.Data = lockMoney
	c4c.SendMsg2Center(baseData)
}

//解锁
func (c4c *Conn4Center) UnlockSettlement(p *Player) {
	timeStr := time.Now().Format("2006-01-02_15:04:05")
	loseOrder := p.Id + "_" + timeStr + "_UnlockMoney"

	baseData := &BaseMessage{}
	baseData.Event = msgUnlockSettlement
	id, _ := strconv.Atoi(p.Id)
	lockMoney := &UserChangeScore{}
	lockMoney.Auth.DevName = conf.Server.DevName
	lockMoney.Auth.DevKey = c4c.DevKey
	lockMoney.Info.CreateTime = time.Now().Unix()
	lockMoney.Info.GameId = c4c.GameId
	lockMoney.Info.ID = id
	lockMoney.Info.LockMoney = p.Account
	lockMoney.Info.Money = 0
	lockMoney.Info.Order = loseOrder
	lockMoney.Info.PayReason = "UnlockMoney"
	lockMoney.Info.PreMoney = 0
	lockMoney.Info.RoundId = p.RoundId
	baseData.Data = lockMoney
	c4c.SendMsg2Center(baseData)
}

func (c4c *Conn4Center) NoticeWinMoreThan(playerId, playerName string, winGold float64) {
	log.Debug("<-------- NoticeWinMoreThan  -------->")
	msg := fmt.Sprintf("<size=20><color=yellow>恭喜!</color><color=orange>%v</color><color=yellow>在</color></><color=orange><size=25>发财推币机</color></><color=yellow><size=20>中一把赢了</color></><color=yellow><size=30>%.2f</color></><color=yellow><size=25>金币！</color></>", playerName, winGold)

	base := &BaseMessage{}
	base.Event = msgWinMoreThanNotice
	id, _ := strconv.Atoi(playerId)
	base.Data = &Notice{
		DevName: conf.Server.DevName,
		DevKey:  conf.Server.DevKey,
		ID:      id,
		GameId:  c4c.GameId,
		Type:    2000,
		Message: msg,
		Topic:   "系统提示",
	}
	c4c.SendMsg2Center(base)
}
