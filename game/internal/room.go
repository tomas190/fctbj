package internal

import (
	"fctbj/msg"
	"fmt"
	"math/rand"
	"time"
)

const (
	GOLD  = 1   // 接金币
	RICH  = 2   // 吐钱
	PUSH  = 3   // 财神推金币
	LUCKY = 4   // 三只小猪
)

var (
	packageTax map[uint16]float64
)

type Room struct {
	RoomId string  // 房间号
	Config string  // 房间配置
	Player *Player // 玩家信息
}

func (r *Room) Init() {
	r.RoomId = fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	r.Config = "1"
}

//RespRoomData 返回房间数据
func (r *Room) RespRoomData() *msg.RoomData {
	rd := &msg.RoomData{}
	rd.RoomId = r.RoomId
	return rd
}
