package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"

	"fctbj/msg"
)

// 定義flag參數，這邊會返回一個相應的指針
var addr = flag.String("addr", "47.56.69.120:1362", "http service address")

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	t := time.NewTicker(time.Millisecond)
	for i := 0; i < 10000; i++ {
		<-t.C
		fmt.Println(i)
		go clientbot(ctx)
	}

	time.Sleep(9000 * time.Second) //等待15分鐘關閉程式

}

func FirstLogin(ws *websocket.Conn) {
	var pkgID uint16
	pkgID = uint16(msg.MessageID_MSG_Login_C2S)
	buf := make([]byte, 100)
	binary.BigEndian.PutUint16(buf[0:2], pkgID)
	playerId := fmt.Sprintf("%08v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(100000000))
	data := msg.Login_C2S{
		Id:       playerId,
		PassWord: "123456",
	}

	// 將資料編碼成 Protocol Buffer 格式（請注意是傳入 Pointer）。
	dataBuffer, _ := proto.Marshal(&data)

	// 將消息ID與DATA整合，一起送出
	pkgData := [][]byte{buf[:2], dataBuffer}
	pkgDatas := bytes.Join(pkgData, []byte{})
	err := ws.WriteMessage(websocket.BinaryMessage, pkgDatas)

	if err != nil {
		log.Println("write:", err)
		return
	}
	// logger.Debug("write:", pkgDatas)
}

func clientbot(ctx context.Context) {

	// 心跳計時器
	heartbeat := time.NewTicker(5 * time.Second)

	// spin計時器
	delayTime := 5
	replay := time.NewTicker(time.Duration(delayTime) * time.Second)

	// 登入後建立的虛擬id
	var userID string

	buf := make([]byte, 100)

	slotCtx, cancel := context.WithCancel(ctx)

	// 處理連接的網址
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	// logger.Debug("connecting to %s", u.String())
	_ = u
	// 連接服務器 本機
	//ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	ws, _, err := websocket.DefaultDialer.Dial("ws://game.539316.com/fctbj", nil)
	// ws, _, err := websocket.DefaultDialer.Dial("ws://game.tampk.club/fctbj", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	defer func() {
		heartbeat.Stop()
		replay.Stop()
		cancel()
		ws.Close() // 預先關閉，此行在離開main時會執行
	}()

	go func(slotCtx context.Context) {
		// 預先關閉，此行在離開本協程時執行
		for {
			select {
			case <-slotCtx.Done():
				return
			default:
				for {
					// 一直待命讀資料
					_, message, err := ws.ReadMessage()
					if err != nil {
						log.Println("read:", err)
						return
					}
					var pkgID = binary.BigEndian.Uint16(message[0:2])
					//var bodyId = binary.BigEndian.Uint16(string(myMsg.Command_Pong))
					if int16(pkgID) == int16(msg.MessageID_MSG_Pong) { // 心跳回傳
						// logger.Debug("MessageKind_Pong")
						// 將已經編碼的資料解碼成 protobuf.User 格式。
						// var bodyClass myMsg.StoCHeartBeat
						// proto.Unmarshal(message[2:], &bodyClass)
						// logger.Debug("心跳回傳recv: %v %v %v %v", pkgID, int16(myMsg.MessageKind_Pong), message, &bodyClass)
					} else if int16(pkgID) == int16(msg.MessageID_MSG_Login_S2C) {
						var bodyClass msg.Login_S2C
						proto.Unmarshal(message[2:], &bodyClass)
						userID = bodyClass.PlayerInfo.Id
						// logger.Debug("登陸回傳recv: %v %v %v %v", pkgID, int16(myMsg.MessageKind_LoginR), message, userID)
					} else {
						// logger.Debug("recv: %v ", message)
					}
				}

			}

		}
	}(slotCtx)

	FirstLogin(ws)
	JoinRoom(ws)
	time.Sleep(1) //確保玩家登入再傳心跳

	// 計時器定時發任務
	for {
		select {
		case <-replay.C:
			//重覆玩
			var pkgID uint16
			pkgID = uint16(msg.MessageID_MSG_PlayerAction_C2S)
			// userIDIndex := rand.Intn(len(allUserID))
			binary.BigEndian.PutUint16(buf[0:2], pkgID)
			mssage := &msg.PlayerAction_C2S{
				DownBet: 0.1,
			}

			// 將資料編碼成 Protocol Buffer 格式（請注意是傳入 Pointer）。
			dataBuffer, _ := proto.Marshal(mssage)

			// 將消息ID與DATA整合，一起送出
			pkgData := [][]byte{buf[:2], dataBuffer}
			pkgDatas := bytes.Join(pkgData, []byte{})
			err = ws.WriteMessage(websocket.BinaryMessage, pkgDatas)

			if err != nil {
				log.Println("write:", err)
				return
			}

		case <-heartbeat.C:
			var pkgID uint16
			pkgID = uint16(msg.MessageID_MSG_Ping)
			binary.BigEndian.PutUint16(buf[0:2], pkgID)
			data := msg.Ping{}

			// 將資料編碼成 Protocol Buffer 格式（請注意是傳入 Pointer）。
			dataBuffer, _ := proto.Marshal(&data)

			// 將消息ID與DATA整合，一起送出
			pkgData := [][]byte{buf[:2], dataBuffer}
			pkgDatas := bytes.Join(pkgData, []byte{})
			err = ws.WriteMessage(websocket.BinaryMessage, pkgDatas)

			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-ctx.Done():
			return

		}
	}
}

func JoinRoom(ws *websocket.Conn) {
	var pkgID uint16
	pkgID = uint16(msg.MessageID_MSG_JoinRoom_C2S)
	buf := make([]byte, 100)
	binary.BigEndian.PutUint16(buf[0:2], pkgID)
	data := msg.JoinRoom_C2S{
		Cfg: "1",
	}

	// 將資料編碼成 Protocol Buffer 格式（請注意是傳入 Pointer）。
	dataBuffer, _ := proto.Marshal(&data)

	// 將消息ID與DATA整合，一起送出
	pkgData := [][]byte{buf[:2], dataBuffer}
	pkgDatas := bytes.Join(pkgData, []byte{})
	err := ws.WriteMessage(websocket.BinaryMessage, pkgDatas)

	if err != nil {
		log.Println("write:", err)
		return
	}
}

func RandInRange(min int, max int) int {
	time.Sleep(1 * time.Nanosecond)
	return rand.Intn(max-min) + min
}
